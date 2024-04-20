package models

import (
	"github.com/Gornak40/crosspawn/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.SqlitePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, err
}
