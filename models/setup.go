package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const sqlitePath = "test.db"

func ConnectDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, err
}
