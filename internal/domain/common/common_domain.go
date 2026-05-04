package common

import "time"

type Pagination struct {
	Page     int   `json:"page" default:"1"`
	PageSize int   `json:"page_size" default:"12"`
	Total    int64 `json:"total"`
}

type Base struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
	CreatedBy uint       `json:"created_by"`
	UpdatedBy uint       `json:"updated_by"`
}
