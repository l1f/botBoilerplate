package modules

type PermissionClass string

const (
	PRIVATE        PermissionClass = "private"
	ADMINISTRATION PermissionClass = "administration"
	MANAGEMENT     PermissionClass = "management"
	GENERAL        PermissionClass = "general"
	FUN            PermissionClass = "fun"

	// this one should never be checked
	PUBLIC PermissionClass = "public"
)

var PermissionClasses = []PermissionClass{
	PRIVATE, ADMINISTRATION, MANAGEMENT, GENERAL, FUN,
}

var restrictedPermissionClasses = []PermissionClass{
	PRIVATE, ADMINISTRATION, MANAGEMENT,
}

func PermissionClassIsRestricted(class PermissionClass) bool {
	for _, permission := range restrictedPermissionClasses {
		if permission == class {
			return true
		}
	}

	return false
}

func PermissionClassExists(class string) bool {
	for _, permission := range PermissionClasses {
		if string(permission) == class {
			return true
		}
	}

	return false
}
