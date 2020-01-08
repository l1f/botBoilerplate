package main

import (
	"botBoilerplate/chat"
	"botBoilerplate/config"
	"botBoilerplate/database"
	"botBoilerplate/messages"
	"botBoilerplate/modules"
	"errors"
	"flag"
	"fmt"
	"log"
	"maunium.net/go/mautrix"
	"os"
	"strings"
	"time"
)

var configPath *string

func init() {
	handleFlags()
}

func main() {
	defer database.Close()

	log.Printf("Logging in to %v as %v", config.Config.Matrix.Homeserver, config.Config.Matrix.Username)

	client, err := connectToHomeServer()
	if err != nil {
		log.Fatalf("Couldn't connect to home server: %v", err)
	}
	log.Println("Login successful")

	modules.Run(client)

	syncer := client.Syncer.(*mautrix.DefaultSyncer)

	setMessageHandler(syncer, client)
	setJoinHandler(syncer, client)

	err = client.Sync()
	if err != nil {
		log.Fatal(err)
	}

}

func handleFlags() {
	help := flag.Bool("h", false, "show help message")
	configPath = flag.String("c", "config.json", "The path to the config file.")
	flag.Parse()

	if *help {
		fmt.Println("A [Matrix]-Bot")
		fmt.Println()
		fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()

		os.Exit(0)
	}
}

func loadConfig() error {
	err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Error while loading config: %v", err)
		return err
	}
	log.Println("Loading config successful")
	if len(config.Config.Matrix.Username) == 0 {
		return errors.New("There is no username in config")
	}
	if len(config.Config.Matrix.Homeserver) == 0 {
		return errors.New("There is no homeserver in config")
	}
	if len(config.Config.Matrix.Password) == 0 {
		return errors.New("There is no password in config")
	}
	return nil
}

func isServerAllowed(roomID string) bool {
	for _, server := range config.Config.AllowedServers {
		if strings.Split(roomID, ":")[1] == server.Homeserver {
			return true
		}
	}
	return true
}

func connectToHomeServer() (*mautrix.Client, error) {
	client, err := mautrix.NewClient(config.Config.Matrix.Homeserver, "", "")
	if err != nil {
		return nil, err
	}
	if len(config.Config.Matrix.AuthToken) == 0 {
		log.Println("No Auth-Token in config")
	}
	resp, err := client.Login(&mautrix.ReqLogin{
		Type:     "m.login.password",
		Password: config.Config.Matrix.Password,
		User:     config.Config.Matrix.Username,
	})
	if err != nil {
		return nil, err
	}
	client.SetCredentials(resp.UserID, resp.AccessToken)

	return client, nil
}

func setJoinHandler(syncer *mautrix.DefaultSyncer, client *mautrix.Client) {
	syncer.OnEventType(mautrix.StateJoinRules, func(event *mautrix.Event) {
		rooms, _ := client.JoinedRooms()
		joined := false
		for _, room := range rooms.JoinedRooms {
			if event.RoomID == room {
				joined = true
				break
			}
		}
		if !joined && isServerAllowed(event.RoomID) {
			_, err := client.JoinRoom(event.RoomID, "", "")
			if err != nil {
				log.Fatalf("Cant join room %v %v\n", event.RoomID, err)
			}
			log.Printf("New room joint, to room %v invited by %v \n", event.RoomID, event.Sender)
			client.SendText(event.RoomID, messages.WELCOME_MESSAGE_1)
			time.Sleep(3 * time.Second)
			client.SendText(event.RoomID, messages.WELCOME_MESSAGE_2)
		}
		if !isServerAllowed(event.RoomID) {
			client.ForgetRoom(event.RoomID)
		}
	})
}

func setMessageHandler(syncer *mautrix.DefaultSyncer, client *mautrix.Client) {
	syncer.OnEventType(mautrix.EventMessage, func(matrixEvent *mautrix.Event) {
		err := client.MarkRead(matrixEvent.RoomID, matrixEvent.ID)
		if err != nil {
			log.Printf("Can't mark message as read: %v", matrixEvent.ID)
		}

		if matrixEvent.Sender == client.UserID {
			// Ignore messages from my self
			return
		}

		go chat.HandleMessage(*client, matrixEvent)
	})
}
