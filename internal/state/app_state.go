package state

import (
	"github.com/redis/go-redis/v9"
	"github.com/starks97/alcohol-tracker-api/config"
	"gorm.io/gorm"
)

type AppState struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
}
