package database

import (
	"fmt"
	"log"

	"github.com/starks97/alcohol-tracker-api/config"
	"github.com/starks97/alcohol-tracker-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.DatabaseUrl
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=UTC",
			cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.PostgresPort,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("✅ Database connected successfully")

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Database connected & migrated successfully")

	return db

}
