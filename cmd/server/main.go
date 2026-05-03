package main

import (
	"go_monolith_sample/internal/pharmacy/delivery"
	"go_monolith_sample/internal/pharmacy/repository"
	"go_monolith_sample/internal/pharmacy/service"
	"go_monolith_sample/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database, _ := db.InitPostgres() // Giả sử em đã bỏ qua check error ở đây cho nhanh

	repo := repository.NewMedicineRepository(database)
	sev := service.NewMedicineService(repo)
	handler := delivery.NewMedicineHandler(sev)

	router := gin.Default()

	api := router.Group("api/v1/medicines")
	{
		api.POST("/", handler.CreateMedicine)
		api.PUT("/:id", handler.UpdateMedicine)
		api.DELETE("/:id", handler.DeleteMedicine)
		api.GET("/:id", handler.GetMedicineByID)
		api.GET("/", handler.GetAllMedicines)
	}

	// 4. Chạy Server
	log.Println("Server đang chạy tại cổng :8080")
	router.Run(":8080")
}
