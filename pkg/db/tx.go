package db

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// Nhét Transaction vào Context
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// Lấy Transaction ra từ Context, nếu không có thì trả về DB mặc định
func GetTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
