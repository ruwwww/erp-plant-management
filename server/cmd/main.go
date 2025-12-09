package main

import (
	"log"
	"os"
	"server/cmd/database"
	v1 "server/http/v1"
	"server/http/v1/handlers"
	"server/internal/core/domain"
	"server/internal/repository"
	"server/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	// Repositories
	userRepo := repository.NewUserRepository(database.DB)
	addrRepo := repository.NewGormRepository[domain.Address](database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	categoryRepo := repository.NewCategoryRepository(database.DB)
	variantRepo := repository.NewVariantRepository(database.DB)
	stockRepo := repository.NewGormRepository[domain.Stock](database.DB)
	movementRepo := repository.NewGormRepository[domain.StockMovement](database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)
	sessionRepo := repository.NewGormRepository[domain.POSSession](database.DB)
	cashMoveRepo := repository.NewGormRepository[domain.POSCashMove](database.DB)
	poRepo := repository.NewGormRepository[domain.PurchaseOrder](database.DB)
	supplierRepo := repository.NewSupplierRepository(database.DB)

	// Services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	userService := service.NewUserService(userRepo, addrRepo)
	catalogService := service.NewCatalogService(productRepo, categoryRepo, variantRepo, database.DB)
	cartService := service.NewCartService()
	inventoryService := service.NewInventoryService(stockRepo, movementRepo, database.DB)
	orderService := service.NewOrderService(orderRepo, inventoryService, database.DB)
	posService := service.NewPOSService(sessionRepo, cashMoveRepo)
	assemblyService := service.NewAssemblyService()
	procurementService := service.NewProcurementService(poRepo, supplierRepo)
	financeService := service.NewFinanceService()

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	storeHandler := handlers.NewStoreHandler(catalogService, cartService, orderService)
	userHandler := handlers.NewUserHandler(userService, orderService, financeService)
	posHandler := handlers.NewPOSHandler(posService, orderService)
	opsHandler := handlers.NewOpsHandler(inventoryService, assemblyService, procurementService)
	adminHandler := handlers.NewAdminHandler(catalogService, authService, userService, procurementService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		ExposeHeaders:    "*",
		AllowCredentials: false,
	}))

	module := os.Getenv("APP_MODULE")
	log.Printf("Starting application with module: %s", module)

	v1.SetupRoutes(app, module, authHandler, storeHandler, userHandler, posHandler, opsHandler, adminHandler)
	log.Fatal(app.Listen(":8080"))
}
