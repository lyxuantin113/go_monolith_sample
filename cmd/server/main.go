package main

import (
	invRepo "go_monolith_sample/internal/inventory/repository"
	medDel "go_monolith_sample/internal/medicine/delivery"
	medRepo "go_monolith_sample/internal/medicine/repository"
	medSer "go_monolith_sample/internal/medicine/service"
	purchDel "go_monolith_sample/internal/purchase/delivery"
	purchRepo "go_monolith_sample/internal/purchase/repository"
	purchSer "go_monolith_sample/internal/purchase/service"
	salesDel "go_monolith_sample/internal/sales/delivery"
	salesRepo "go_monolith_sample/internal/sales/repository"
	salesSer "go_monolith_sample/internal/sales/service"

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
	invRepo := invRepo.NewInventoryRepository(database)

	medRepo := medRepo.NewMedicineRepository(database)
	medSer := medSer.NewMedicineService(medRepo)
	medHan := medDel.NewMedicineHandler(medSer)

	salesRepo := salesRepo.NewOrderRepository(database)
	salesSer := salesSer.NewOrderService(salesRepo, medSer, invRepo)
	salesHan := salesDel.NewOrderHandler(salesSer)

	purchRepo := purchRepo.NewPurchaseOrderRepository(database)
	purchSer := purchSer.NewPurchaseOrderService(purchRepo, medSer, invRepo)
	purchHan := purchDel.NewPurchaseOrderHandler(purchSer)

	// 2. Khởi tạo Gin
	r := gin.Default()

	// 3. Đăng ký Routes
	api := r.Group("/api/v1")
	{
		medicines := api.Group("/medicines")
		{
			medicines.POST("", medHan.CreateMedicine)
			medicines.GET("", medHan.GetAllMedicines)
			medicines.GET("/:id", medHan.GetMedicineByID)
			medicines.PUT("/:id", medHan.UpdateMedicine)
			medicines.DELETE("/:id", medHan.DeleteMedicine)
		}

		sales := api.Group("/sales")
		{
			sales.POST("", salesHan.CreateOrder)
			sales.PUT("/:id/refund", salesHan.RefundOrder)
		}

		purchases := api.Group("/purchases")
		{
			purchases.POST("", purchHan.CreatePurchaseOrder)
			purchases.PUT("/:id/complete", purchHan.CompletePurchaseOrder)
			purchases.PUT("/:id/cancel", purchHan.CancelPurchaseOrder)
		}
	}

	// 4. Chạy Server
	log.Println("Server đang chạy tại port 8080...")
	r.Run(":8080")
}
