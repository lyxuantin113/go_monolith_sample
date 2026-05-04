package auth

import (
	"go_monolith_sample/internal/domain/common"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleStaff    UserRole = "staff"
	RoleCustomer UserRole = "customer"
)

type User struct {
	common.Base
	Username string   `gorm:"uniqueIndex;notnull" json:"username"`
	Password string   `json:"-"` // Luôn thêm json:"-" để không lộ mật khẩu khi trả về API
	FullName string   `json:"full_name"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Address  string   `json:"address"`
	Role     UserRole `gorm:"type:varchar(20);default:'customer'" json:"role"`
}

// Nếu là Khách hàng thì có thêm thông tin tích điểm
type CustomerProfile struct {
	common.Base
	UserID   uint   `gorm:"uniqueIndex" json:"user_id"`
	Points   int    `gorm:"default:0" json:"points"`
	FullName string `json:"full_name"`
}
