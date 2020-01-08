package commands

import (
	"botBoilerplate/modules"
	"botBoilerplate/modules/options"
	"botBoilerplate/modules/permissions"
	"fmt"
	"maunium.net/go/mautrix"
	"strings"
)

const INTERNAL_HELP_TEXT = `Hi there!
I am a Matrix Bot.

Here are some commands I can execute in this room:
%s
`

func init() {
	modules.RegisterCommand(modules.Command{
		Name:        "ping",
		Description: `ping-pong`,
		Class:       modules.FUN,
		Handler:     ping,
	})
	modules.RegisterCommand(modules.Command{
		Name:        "help",
		Description: `displays help`,
		Class:       modules.GENERAL,
		Handler:     help,
	})
}

func ping(client mautrix.Client, event *mautrix.Event, options *options.CommandOptions) modules.CommandResult {
	client.SendText(event.RoomID, "Pong")

	return modules.SUCCESS
}

func help(client mautrix.Client, event *mautrix.Event, options *options.CommandOptions) modules.CommandResult {
	var builder strings.Builder

	commands := modules.GetCommands()

	permissions := permissions.Get(event.RoomID)

	for _, command := range commands {
		if command.Class != modules.PUBLIC {
			ok := false
			for _, permission := range permissions {
				if permission == command.Class {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}

		builder.WriteString(" - ")
		builder.WriteString(command.Name)
		builder.WriteString(" (")
		builder.WriteString(string(command.Class))
		builder.WriteString("): ")
		builder.WriteString(command.Description)
		builder.WriteString("\n")
	}

	client.SendText(event.RoomID, fmt.Sprintf(INTERNAL_HELP_TEXT, builder.String()))

	return modules.SUCCESS
}