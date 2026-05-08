package main

import (
	"go_monolith_sample/internal/domain/common"
	invRepo "go_monolith_sample/internal/inventory/repository"
	medDel "go_monolith_sample/internal/medicine/delivery"
	medRepo "go_monolith_sample/internal/medicine/repository"
	medSer "go_monolith_sample/internal/medicine/service"
	purchDel "go_monolith_sample/internal/purchase/delivery"
	purchRepo "go_monolith_sample/internal/purchase/repository"
	purchSer "go_monolith_sample/internal/purchase/service"
	reportDel "go_monolith_sample/internal/report/delivery"
	reportSer "go_monolith_sample/internal/report/service"
	salesDel "go_monolith_sample/internal/sales/delivery"
	salesRepo "go_monolith_sample/internal/sales/repository"
	salesSer "go_monolith_sample/internal/sales/service"

	reportRepo "go_monolith_sample/internal/report/repository"
	authMiddleware "go_monolith_sample/pkg/auth/middleware"
	"go_monolith_sample/pkg/cache/redis"

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

	reportRepo := reportRepo.NewReportRepository(database)
	reportSer := reportSer.NewReportService(medSer, reportRepo)
	reportHan := reportDel.NewReportHandler(reportSer)

	// 2. Khởi tạo Gin
	r := gin.Default()
	redisClient := redis.NewRedisClient("localhost:6379")

	// 3. Đăng ký Routes
	api := r.Group("/api/v1")
	{
		adminRoutes := api.Group("/")
		adminRoutes.Use(authMiddleware.Authorize(common.RoleAdmin))

		managerRoutes := api.Group("/")
		managerRoutes.Use(authMiddleware.Authorize(common.RoleAdmin, common.RoleManager))

		staffRoutes := api.Group("/")
		staffRoutes.Use(authMiddleware.Authorize(common.RoleAdmin, common.RoleManager, common.RoleStaff))

		medicines := api.Group("/medicines")
		{
			medicines.POST("", medHan.CreateMedicine)
			medicines.GET("", medHan.GetAllMedicines)
			medicines.GET("/:id", medHan.GetMedicineByID)
			medicines.PUT("/:id", medHan.UpdateMedicine)

			// Admin
			medicinesAdmin := medicines.Group("/")
			medicinesAdmin.Use(authMiddleware.Authorize(common.RoleAdmin))
			medicinesAdmin.DELETE("/:id", medHan.DeleteMedicine)
		}

		sales := api.Group("/sales")
		{
			salesAdmin := sales.Group("/")
			salesAdmin.Use(authMiddleware.Authorize(common.RoleAdmin))
			salesManager := sales.Group("/")
			salesManager.Use(authMiddleware.Authorize(common.RoleAdmin, common.RoleManager))

			sales.POST("", salesHan.CreateOrder)

			// Manager
			salesManager.PUT("/:id/refund", salesHan.RefundOrder)

			// Admin
			salesAdmin.DELETE("/:id", salesHan.DeleteOrder)
		}

		purchases := api.Group("/purchases")
		{
			// Admin
			purchasesManager := purchases.Group("/")
			purchasesManager.Use(authMiddleware.Authorize(common.RoleAdmin, common.RoleManager))
			purchasesManager.POST("", purchHan.CreatePurchaseOrder)
			purchasesManager.PUT("/:id/complete", purchHan.CompletePurchaseOrder)
			purchasesManager.PUT("/:id/cancel", purchHan.CancelPurchaseOrder)
		}

		// REPORT
		reports := api.Group("/reports")
		reports.Use(authMiddleware.Authenticate(redisClient), authMiddleware.Authorize(common.RoleAdmin, common.RoleManager))
		{
			reports.GET("/ledger", reportHan.GetInventoryLedger)
			reports.GET("/sio", reportHan.GetSIOReport)
			reports.GET("/sio/export", reportHan.ExportSIOExcel)
			reports.POST("/snapshot", reportHan.TakeSnapshot)
		}

	}

	// 4. Chạy Server
	log.Println("Server đang chạy tại port 8080...")
	r.Run(":8080")
}
