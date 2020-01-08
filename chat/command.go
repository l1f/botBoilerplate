package chat

import (
	"botBoilerplate/messages"
	"botBoilerplate/modules"
	"botBoilerplate/modules/options"
	"botBoilerplate/modules/permissions"
	"fmt"
	"log"
	"maunium.net/go/mautrix"
	"strings"
)

func parseCommand(client mautrix.Client, event *mautrix.Event) {
	text := strings.TrimPrefix(event.Content.Body, "!")

	options, err := options.New(text)
	if err != nil {
		client.SendText(event.RoomID, err.Error())
		return
	}

	command, err := modules.FindCommand(options.Name)
	if err != nil {
		client.SendText(event.RoomID, err.Error())
		return
	}

	if command.Class != modules.PUBLIC && !permissions.Check(event.RoomID, command.Class) {
		log.Printf("Command %s not allowed in room %s.", command.Name, event.RoomID)
		client.SendText(event.RoomID, fmt.Sprintf(messages.COMMAND_NOT_ALLOWED))
		return
	}

	log.Printf("Executing command %s in room %s", command.Name, event.RoomID)

	result := command.Handler(client, event, options)
	switch result {
	case modules.SUCCESS:
		log.Printf("Command terminated successfully.")
	case modules.FAIL:
		log.Printf("Command failed.")
		client.SendText(event.RoomID, fmt.Sprintf(messages.COMMAND_FAILED, command.Name))
	case modules.SILENT_FAIL:
		log.Printf("Command failed silently.")
	default:
		log.Printf("Unknown command result")
	}
}

func isCommand(event *mautrix.Event) bool {
	return strings.HasPrefix(event.Content.Body, "!")
}
