package repository

import (
	"time"

	"github.com/namchee/go-oauth2-client/models"
	"golang.org/x/oauth2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db gorm.DB
)

func init() {
	dbInstance, err := gorm.Open(sqlite.Open("../oauth2.sql"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db = *dbInstance

	db.AutoMigrate(&models.Token{})
}

func GetToken(sessionToken string) (models.Token, error) {
	var tokenMap models.Token

	if err := db.Where("session_token = ?", sessionToken).Find(&tokenMap).Error; err != nil {
		return tokenMap, err
	}

	return tokenMap, nil
}

func CreateNewToken(
	sessionToken string,
	token *oauth2.Token,
	updatedAt *time.Time,
) (models.Token, error) {
	tokenMap := models.Token{
		SessionToken: sessionToken,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TimeToLive:   &token.Expiry,
		UpdatedAt:    updatedAt,
	}

	result := db.Create(&tokenMap)

	if result.Error != nil {
		return tokenMap, result.Error
	}

	return tokenMap, nil
}

func RefreshToken(
	sessionToken string,
	accessToken string,
	updatedAt *time.Time,
) (models.Token, error) {
	token, err := GetToken(sessionToken)

	if err != nil {
		return token, err
	}

	db.Model(&token).Update("access_token", accessToken)

	return token, nil
}
