package modules

import (
	"botBoilerplate/messages"
	"botBoilerplate/modules/options"
	"errors"
	"maunium.net/go/mautrix"
)

type CommandResult int

const (
	SUCCESS     CommandResult = iota
	FAIL        CommandResult = iota
	SILENT_FAIL CommandResult = iota
)

type CommandHandler func(client mautrix.Client, event *mautrix.Event, options *options.CommandOptions) CommandResult

type Command struct {
	Name        string
	Description string
	Class       PermissionClass
	Handler     CommandHandler
}

var commands = map[string]Command{}

func RegisterCommand(command Command) {
	commands[command.Name] = command
}

func FindCommand(name string) (Command, error) {
	command, okay := commands[name]
	if !okay {
		return command, errors.New(messages.COMMAND_NOT_FOUND)
	}

	return command, nil
}

func GetCommands() []Command {
	var commandList []Command

	for _, command := range commands {
		commandList = append(commandList, command)
	}

	return commandList
}
