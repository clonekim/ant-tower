package data

import (
	"fmt"
	"time"

	"twn-monitor/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ProcessLog struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	PID       int32      `json:"pid"`
	Name      string     `json:"name"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Duration  string     `json:"duration"`
}

var DB *gorm.DB

func InitDB(cfg *config.Config) error {

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBFile), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to DB (%s): %w", cfg.DBFile, err)
	}

	if err := DB.AutoMigrate(&ProcessLog{}); err != nil {
		return fmt.Errorf("failed to migrate DB schema: %w", err)
	}

	log.Info().Str("path", cfg.DBFile).Msg("DB initialized")
	return nil
}
