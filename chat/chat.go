package chat

import (
	"botBoilerplate/modules"
	"maunium.net/go/mautrix"
)

func HandleMessage(client mautrix.Client, matrixEvent *mautrix.Event) {
	if isCommand(matrixEvent) {
		parseCommand(client, matrixEvent)
	}

	modules.FindBuzzword(client, matrixEvent)
}