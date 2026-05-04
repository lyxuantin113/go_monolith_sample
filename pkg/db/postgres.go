package db

import (
	"fmt"
	"os"

	"go_monolith_sample/internal/domain/inventory"
	med "go_monolith_sample/internal/domain/medicine"
	sales "go_monolith_sample/internal/domain/sales"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() (*gorm.DB, error) {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&med.Medicine{}, &sales.Order{}, &sales.OrderItem{}, &inventory.InventoryTransaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
