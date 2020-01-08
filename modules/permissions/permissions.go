package permissions

import (
	"botBoilerplate/config"
	"botBoilerplate/database"
	"botBoilerplate/modules"
	"github.com/jinzhu/gorm"
)

type RoomPermission struct {
	gorm.Model
	Room string
	Permission modules.PermissionClass
	Enabled bool
	ChangedBy string
}

func init() {
	modules.RegisterModel(modules.Model{
		Name:  "RoomPermissions",
		Model: &RoomPermission{},
	})
}

func Get(room string) []modules.PermissionClass {
	db := database.GetDB()

	var roomPermissions []RoomPermission

	var permissionsMap = map[modules.PermissionClass]bool{}
	
	for _, permission := range config.Config.DefaultPermissions {
		permissionsMap[permission] = true
	}
	
	db.Where(&RoomPermission{
		Room:   room,
	}).Find(&roomPermissions)

	for _, roomPermission := range roomPermissions {
		if roomPermission.Enabled {
			permissionsMap[roomPermission.Permission] = true
		} else {
			delete(permissionsMap, roomPermission.Permission)
		}
	}

	var permissions []modules.PermissionClass

	for permission, _ := range permissionsMap {
		permissions = append(permissions, permission)
	}

	return permissions
}

func Check(room string, class modules.PermissionClass) bool {
	db := database.GetDB()

	var permissions []RoomPermission

	db.Where(&RoomPermission{Room:room, Permission:class}).Find(&permissions)

	if len(permissions) >= 1 {
		return permissions[0].Enabled
	}

	for _, permission := range config.Config.DefaultPermissions {
		if permission == class {
			return true
		}
	}

	return false
}

func IsPrivileged(user string) bool {
	for _, privileged := range config.Config.PrivilegedUsers {
		if privileged == user {
			return true
		}
	}

	return false
}

func Add(room string, user string, class modules.PermissionClass) {
	db := database.GetDB()

	var roomPermissions []RoomPermission

	db.Where(&RoomPermission{Room: room, Permission: class}).Find(&roomPermissions)

	if len(roomPermissions) > 0 {
		roomPermissions[0].Enabled = true
		roomPermissions[0].ChangedBy = user

		db.Save(&(roomPermissions[0]))
	} else {
		roomPermission := RoomPermission{
			Room:       room,
			Permission: class,
			Enabled:    true,
			ChangedBy:  user,
		}

		db.Create(&roomPermission)
	}
}

func Delete(room string, user string, class modules.PermissionClass) {
	db := database.GetDB()

	var roomPermissions []RoomPermission

	db.Where(&RoomPermission{Room: room, Permission: class}).Find(&roomPermissions)

	if len(roomPermissions) > 0 {
		roomPermissions[0].Enabled = false
		roomPermissions[0].ChangedBy = user

		db.Save(&(roomPermissions[0]))
	} else {
		roomPermission := RoomPermission{
			Room:       room,
			Permission: class,
			Enabled:    false,
			ChangedBy:  user,
		}

		db.Create(&roomPermission)
	}
}
