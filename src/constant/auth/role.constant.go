package auth

type Role string

const (
	ADMIN       Role = "admin"
	EVENT_STAFF      = "event_staff"
	BAAN_STAFF       = "baan_staff"
	USER             = "user"
)
