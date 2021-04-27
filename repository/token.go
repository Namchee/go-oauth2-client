package repository

import (
	"time"

	"github.com/namchee/go-oauth2-client/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db gorm.DB
)

func init() {
	dbInstance, err := gorm.Open(sqlite.Open("oauth2.sql"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db = *dbInstance

	db.AutoMigrate(&models.TokenMap{})
}

func GetToken(sessionToken string) (models.TokenMap, error) {
	var tokenMap models.TokenMap

	if err := db.Where("session_token = ?", sessionToken).Find(&tokenMap).Error; err != nil {
		return tokenMap, err
	}

	return tokenMap, nil
}

func CreateNewToken(
	sessionToken string,
	accessToken string,
	createdAt *time.Time,
) (models.TokenMap, error) {
	tokenMap := models.TokenMap{
		SessionToken: sessionToken,
		AccessToken:  accessToken,
		CreatedAt:    createdAt,
	}

	result := db.Create(&tokenMap)

	if result.Error != nil {
		return tokenMap, result.Error
	}

	return tokenMap, nil
}
