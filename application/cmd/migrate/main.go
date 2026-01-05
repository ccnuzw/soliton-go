package main

import (
	"fmt"
	"os"

	"gorm.io/gorm"

	"github.com/soliton-go/framework/core/config"
	"github.com/soliton-go/framework/core/logger"
	"github.com/soliton-go/framework/orm"

	userapp "github.com/soliton-go/application/internal/application/user"
	orderapp "github.com/soliton-go/application/internal/application/order"
	productapp "github.com/soliton-go/application/internal/application/product"
	// soliton-gen:imports
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config:", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create logger:", err)
		os.Exit(1)
	}

	db, err := orm.NewGormDB(cfg, log)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect db:", err)
		os.Exit(1)
	}

	if err := migrateAll(db); err != nil {
		fmt.Fprintln(os.Stderr, "migration failed:", err)
		os.Exit(1)
	}

	fmt.Println("migration completed")
}

func migrateAll(db *gorm.DB) error {
	if err := userapp.RegisterMigration(db); err != nil {
		return err
	}
	if err := orderapp.RegisterMigration(db); err != nil {
		return err
	}
	if err := productapp.RegisterMigration(db); err != nil {
		return err
	}
	// soliton-gen:migrations
	return nil
}
