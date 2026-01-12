package orm

import (
	"fmt"

	"github.com/soliton-go/framework/core/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewGormDB creates a new GORM database connection.
func NewGormDB(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	driver := cfg.GetString("database.driver")
	dsn := cfg.GetString("database.dsn")

	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		// Default fallback to sqlite in memory if not configured, or error out
		if driver == "" {
			logger.Info("No database driver specified, defaulting to sqlite in-memory")
			dialector = sqlite.Open("file::memory:?cache=shared")
		} else {
			return nil, fmt.Errorf("unsupported database driver: %s", driver)
		}
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
