package database

import (
	"fmt"
	"log"

	"github.com/starks97/alcohol-tracker-api/config"
	"github.com/starks97/alcohol-tracker-api/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global variable representing the database connection.
// It should be initialized by calling ConnectDB.
var DB *gorm.DB

// ConnectDB establishes a connection to the PostgreSQL database using the provided configuration.
// It also performs automatic database migrations using GORM.
//
// Parameters:
//   - cfg: *config.Config - The application configuration containing database connection details.
//
// Returns:
//   - *gorm.DB: A pointer to the initialized GORM database connection.
//   - It also sets the global var DB.
func ConnectDB(cfg *config.Config) *gorm.DB {
	// Open a connection to the PostgreSQL database using GORM.
	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Log a successful database connection message.
	fmt.Println("✅ Database connected successfully")

	// Perform automatic database migrations for the User model.
	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Log a successful database connection and migration message.
	fmt.Println("✅ Database connected & migrated successfully")

	DB = db

	// Return the initialized GORM database connection.
	return db
}
