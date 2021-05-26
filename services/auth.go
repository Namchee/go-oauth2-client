package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/namchee/go-oauth2-client/models"
	"github.com/namchee/go-oauth2-client/repository"
	"github.com/namchee/go-oauth2-client/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	authConfig *oauth2.Config
	authState  string
)

func init() {
	credentials := utils.GetCredentials()

	authConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/oauth/redirect",
		ClientID:     credentials.ClientId,
		ClientSecret: credentials.ClientSecret,
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}

	randomString, err := uuid.NewRandom()

	if err != nil {
		log.Fatalln("Failed to generate OAuth state")
	}

	authState = randomString.String()
}

func GetLoginUrl() string {
	return authConfig.AuthCodeURL(authState)
}

func ValidateState(state string) bool {
	return state == authState
}

func GetSessionToken(authCode string) (string, error) {
	token, err := authConfig.Exchange(
		context.TODO(),
		authCode,
		oauth2.AccessTypeOffline, // ensures refresh token availability
	)

	if err != nil {
		return "", err
	}

	sessionToken, err := utils.GenerateRandomString()

	if err != nil {
		return "", err
	}

	now := time.Now()

	tokenMap, err := repository.CreateNewToken(sessionToken, token, &now)

	if err != nil {
		return "", err
	}

	return tokenMap.SessionToken, nil
}

func MapToken(sessionToken string) (models.Token, error) {
	tokenMap, err := repository.GetToken(sessionToken)

	if err != nil {
		return tokenMap, err
	}

	oauthToken := oauth2.Token{
		AccessToken:  tokenMap.AccessToken,
		RefreshToken: tokenMap.RefreshToken,
		TokenType:    "Bearer",
	}

	tokenSource := authConfig.TokenSource(context.TODO(), &oauthToken)
	token, err := tokenSource.Token()

	if err != nil {
		return tokenMap, err
	}

	if token != &oauthToken {
		now := time.Now()

		newTokenMap, err := repository.RefreshToken(
			tokenMap.SessionToken,
			token.AccessToken,
			&now,
		)

		if err != nil {
			return newTokenMap, err
		}

		tokenMap = newTokenMap
	}

	return tokenMap, nil
}

func Logout(sessionToken string) error {
	repository.DeleteToken(sessionToken)

	return nil
}
