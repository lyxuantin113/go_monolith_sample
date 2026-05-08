package common

type UserRole string

const (
	RoleAdmin   UserRole = "ADMIN"
	RoleManager UserRole = "MANAGER"
	RoleStaff   UserRole = "STAFF"
)
