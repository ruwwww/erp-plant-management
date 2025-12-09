package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"server/internal/core/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	for i := 0; i < 30; i++ { // Retry up to 30 times
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: false, // Disable prepared statements for pgbouncer compatibility
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to database after retries: ", err)
	}

	log.Println("Database connected successfully")

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.Category{},
		&domain.Product{},
		&domain.ProductVariant{},
		&domain.Tag{},
		&domain.ProductTag{},
		&domain.Stock{},
		&domain.StockMovement{},
		&domain.StockAssembly{},
		&domain.ProductRecipe{},
		&domain.SalesOrder{},
		&domain.SalesOrderItem{},
		&domain.POSSession{},
		&domain.POSCashMove{},
		&domain.PurchaseOrder{},
		&domain.PurchaseOrderItem{},
		&domain.Supplier{},
		&domain.Invoice{},
		&domain.Promotion{},
		&domain.PromotionUsage{},
		&domain.SalesOrderPromotion{},
		&domain.MediaAsset{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")
}
