package state

import (
	"github.com/redis/go-redis/v9"
	"github.com/starks97/alcohol-tracker-api/config"
	"gorm.io/gorm"
)

// AppState holds the application's global state and dependencies, such as the database connection,
// Redis client, and configuration.
type AppState struct {
	DB     *gorm.DB       // Database connection pool.
	Redis  *redis.Client  // Redis client for caching and session management.
	Config *config.Config // Application configuration.
}
