package modules

import (
	"math/rand"
	"maunium.net/go/mautrix"
	"strings"
)

type Buzzword struct {
	Name        string
	Class       PermissionClass
	Trigger     []string
	Reaction    []string
	Probability float32
}

var buzzwords = map[string]Buzzword{}

func RegisterBuzzword(buzzword Buzzword) {
	buzzwords[buzzword.Name] = buzzword
}

func FindBuzzword(client mautrix.Client, event *mautrix.Event) {
	for _, bzzwrd := range buzzwords {
		for _, word := range bzzwrd.Trigger {
			if strings.Contains(event.Content.Body, word) {
				buzzword, _ := buzzwords[bzzwrd.Name]
				sendReaction(client, event, buzzword)
			}
		}
	}
}

func sendReaction(client mautrix.Client, event *mautrix.Event, buzzword Buzzword) {
	if rand.Float32() < buzzword.Probability {
		client.SendText(event.RoomID, buzzword.Reaction[rand.Intn(len(buzzword.Reaction))])
	}
}
