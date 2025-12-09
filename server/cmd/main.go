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
	productRepo := repository.NewGormRepository[domain.Product](database.DB)
	categoryRepo := repository.NewGormRepository[domain.Category](database.DB)
	stockRepo := repository.NewGormRepository[domain.Stock](database.DB)
	movementRepo := repository.NewGormRepository[domain.StockMovement](database.DB)

	// Services
	authService := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	catalogService := service.NewCatalogService(productRepo, categoryRepo)
	inventoryService := service.NewInventoryService(stockRepo, movementRepo, database.DB)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	storeHandler := handlers.NewStoreHandler(catalogService)
	userHandler := handlers.NewUserHandler()
	posHandler := handlers.NewPOSHandler()
	opsHandler := handlers.NewOpsHandler(inventoryService)
	adminHandler := handlers.NewAdminHandler()

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
