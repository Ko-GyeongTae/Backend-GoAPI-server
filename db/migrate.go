package db

import (
	"github.com/Backend-GoAPI-server/model"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.LogMode(false)

	// Create table when table is not exist
	if !db.HasTable(&model.User{}) {
		db.AutoMigrate(&model.User{})
	}
}
