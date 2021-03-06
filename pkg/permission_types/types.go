package permission_types

type Permission string

var (
	StatusNoAccess    Permission = ""
	StatusReadAccess  Permission = "read"
	StatusWriteAccess Permission = "write"
	StatusAdminAccess Permission = "admin"
	StatusOwnerAccess Permission = "owner"
)

func NoAccess() Permission {
	return StatusNoAccess
}

func ReadAccess() Permission {
	return StatusReadAccess
}

func WriteAccess() Permission {
	return StatusWriteAccess
}

func AdminAccess() Permission {
	return StatusAdminAccess
}

func OwnerAccess() Permission {
	return StatusOwnerAccess
}
