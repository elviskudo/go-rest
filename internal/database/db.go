package database

import (
	"go-rest/internal/models"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database, err := gorm.Open(sqlite.Open("inventory.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	// Drop tables to migrate to UUID
	// WARNING: This deletes all data!
	database.Migrator().DropTable(&models.Item{}, &models.User{}, &models.Warehouse{}, &models.Supplier{}, &models.Discount{}, &models.Media{}, &models.Variant{}, &models.Option{}, &models.Review{}, &models.Favorite{}, &models.Inventory{}, &models.Category{}, &models.PurchaseOrder{}, &models.PurchaseOrderItem{}, &models.Order{}, &models.OrderItem{})

	err = database.AutoMigrate(&models.Item{}, &models.User{}, &models.Warehouse{}, &models.Supplier{}, &models.Discount{}, &models.Media{}, &models.Variant{}, &models.Option{}, &models.Review{}, &models.Favorite{}, &models.Inventory{}, &models.Category{}, &models.PurchaseOrder{}, &models.PurchaseOrderItem{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database!", err)
	}

	DB = database
}
