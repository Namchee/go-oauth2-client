package utils

import (
	"github.com/namchee/go-oauth2-client/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDb() (gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("oauth2.sql"), &gorm.Config{})

	if err != nil {
		return *db, err
	}

	return *db, nil
}

func MigrateDb(db gorm.DB) {
	// demo purposes, don't do this
	db.Migrator().DropTable(&models.Token{})
	db.Migrator().CreateTable(&models.Token{})
}
