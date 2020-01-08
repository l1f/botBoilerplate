package commands

import (
	"botBoilerplate/modules"
	"botBoilerplate/modules/options"
	"botBoilerplate/modules/permissions"
	"fmt"
	"maunium.net/go/mautrix"
	"strings"
)

func init() {
	modules.RegisterCommand(modules.Command{
		Name:        "Permission",
		Description: `manages authorization classes`,
		Class:       modules.PUBLIC,
		Handler:     permissionCommand,
	})
}

const PERMISSIONS_HELP_TEXT = `This program manages permissions for this bot.

Usage: !permission [OPTIONS]

Options:

%s

The following authorization classes exist:
(Those marked with * can only be changed by authorized persons)

%s
`

func buildPermissionList(permissions []modules.PermissionClass) string {
	var builder strings.Builder

	for _, permission := range permissions {
		builder.WriteString(" - ")
		builder.WriteString(string(permission))

		if modules.PermissionClassIsRestricted(permission) {
			builder.WriteString("*")
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func permissionCommand(client mautrix.Client, event *mautrix.Event, options *options.CommandOptions) modules.CommandResult {
	help := options.Bool("help", false, "displays this help")

	list := options.Bool("list", false, "Displays a list of the active permission class for this room")

	add := options.String("add", "", "adds a new authorization class to this room")
	delete := options.String("delete", "", "Deletes the specified authorization class for this room")

	err := options.Parse()
	if err != nil {
		client.SendText(event.RoomID, fmt.Sprintf("An error occurred while parsing the arguments: %v", err))
		return modules.SILENT_FAIL
	}

	if *help {
		client.SendText(event.RoomID, fmt.Sprintf(PERMISSIONS_HELP_TEXT, options.Help(), buildPermissionList(modules.PermissionClasses)))
		return modules.SUCCESS
	}

	if *list {
		client.SendText(event.RoomID, fmt.Sprintf("This room has the following authorization classes:\n(Those marked with * can only be changed by authorized persons.)\n\n%s", buildPermissionList(permissions.Get(event.RoomID))))
		return modules.SUCCESS
	}

	if *add == "" && *delete == "" {
		client.SendText(event.RoomID, "To get help with this command, use -help.")
		return modules.SILENT_FAIL
	}

	if *add != "" {
		if !modules.PermissionClassExists(*add) {
			client.SendText(event.RoomID, fmt.Sprintf("The specified authorization class does not exist.\n\nHere is a list of all classes:\n(Those marked with * can only be changed by authorized persons.)\n\n%s", buildPermissionList(modules.PermissionClasses)))
			return modules.SILENT_FAIL
		}

		permission := modules.PermissionClass(*add)

		if permissions.Check(event.RoomID, permission) {
			client.SendText(event.RoomID, "This room already has the specified authorization class.")
			return modules.SILENT_FAIL
		}

		if modules.PermissionClassIsRestricted(permission) && !permissions.IsPrivileged(event.Sender) {
			client.SendText(event.RoomID, "Unfortunately, you are not authorized to change this authorization class.")
			return modules.SILENT_FAIL
		}

		permissions.Add(event.RoomID, event.Sender, permission)

		client.SendText(event.RoomID, fmt.Sprintf("The authorization class %s was added successfully.", permission))
	}

	if *delete != "" {
		if !modules.PermissionClassExists(*delete) {
			client.SendText(event.RoomID, fmt.Sprintf("The specified authorization class does not exist.\n\nHere is a list of all classes:\n(Those marked with * can only be changed by authorized persons.)\n\n%s", buildPermissionList(modules.PermissionClasses)))
			return modules.SILENT_FAIL
		}

		permission := modules.PermissionClass(*delete)

		if !permissions.Check(event.RoomID, permission) {
			client.SendText(event.RoomID, "This room does not have the specified authorization class.")
			return modules.SILENT_FAIL
		}

		if modules.PermissionClassIsRestricted(permission) && !permissions.IsPrivileged(event.Sender) {
			client.SendText(event.RoomID, "Unfortunately, you are not authorized to change this authorization class.")
			return modules.SILENT_FAIL
		}

		permissions.Delete(event.RoomID, event.Sender, permission)

		client.SendText(event.RoomID, fmt.Sprintf("The authorization class %s was successfully removed.", permission))
	}

	return modules.SUCCESS
}
