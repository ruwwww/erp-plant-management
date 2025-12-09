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
	userRepo := repository.NewGormRepository[domain.User](database.DB)
	addrRepo := repository.NewGormRepository[domain.Address](database.DB)
	productRepo := repository.NewGormRepository[domain.Product](database.DB)
	categoryRepo := repository.NewGormRepository[domain.Category](database.DB)
	stockRepo := repository.NewGormRepository[domain.Stock](database.DB)
	movementRepo := repository.NewGormRepository[domain.StockMovement](database.DB)
	orderRepo := repository.NewGormRepository[domain.SalesOrder](database.DB)
	sessionRepo := repository.NewGormRepository[domain.POSSession](database.DB)
	poRepo := repository.NewGormRepository[domain.PurchaseOrder](database.DB)

	// Services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	userService := service.NewUserService(userRepo, addrRepo)
	catalogService := service.NewCatalogService(productRepo, categoryRepo)
	cartService := service.NewCartService()
	orderService := service.NewOrderService(orderRepo)
	inventoryService := service.NewInventoryService(stockRepo, movementRepo, database.DB)
	posService := service.NewPOSService(sessionRepo)
	assemblyService := service.NewAssemblyService()
	procurementService := service.NewProcurementService(poRepo)
	financeService := service.NewFinanceService()

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	storeHandler := handlers.NewStoreHandler(catalogService, cartService, orderService)
	userHandler := handlers.NewUserHandler(userService, orderService, financeService)
	posHandler := handlers.NewPOSHandler(posService)
	opsHandler := handlers.NewOpsHandler(inventoryService, assemblyService, procurementService)
	adminHandler := handlers.NewAdminHandler(catalogService, authService, userService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		ExposeHeaders:    "*",
		AllowCredentials: true,
	}))

	v1.SetupRoutes(app, authHandler, storeHandler, userHandler, posHandler, opsHandler, adminHandler)
	app.Listen("8080")
}
