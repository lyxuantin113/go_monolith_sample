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
	database, err := db.InitPostgres()
	if err != nil {
		log.Fatal(err)
	}

	// 1. Khởi tạo các lớp (Wiring)
	repo := repository.NewMedicineRepository(database)
	sev := service.NewMedicineService(repo)
	handler := delivery.NewMedicineHandler(sev)

	// 2. Khởi tạo Gin
	r := gin.Default()

	// 3. Đăng ký Routes
	api := r.Group("/api/v1")
	{
		medicines := api.Group("/medicines")
		{
			medicines.POST("", handler.CreateMedicine)
			medicines.GET("", handler.GetAllMedicines)
			medicines.GET("/:id", handler.GetMedicineByID)
			medicines.PUT("/:id", handler.UpdateMedicine)
			medicines.DELETE("/:id", handler.DeleteMedicine)
		}
	}

	// 4. Chạy Server
	log.Println("Server đang chạy tại port 8080...")
	r.Run(":8080")
}
