package user

import (
	"context"
	"go_monolith_sample/internal/domain/common"
)

type User struct {
	common.Base
	Username string          `gorm:"uniqueIndex;notnull" json:"username"`
	Password string          `json:"-"`
	FullName string          `json:"full_name"`
	Role     common.UserRole `gorm:"type:varchar(20);default:'STAFF'" json:"role"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
}

type UserService interface {
	Register(ctx context.Context, u *User) error
	Login(ctx context.Context, username, password string) (string, string, error) // Trả về Access & Refresh Token
}
