package services

import (
	"log"
	"time"

	"github.com/google/uuid"
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
		RedirectURL:  "http://localhost:8080/auth/callback",
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
	token, err := authConfig.Exchange(oauth2.NoContext, authCode)

	if err != nil {
		return "", err
	}

	sessionToken, err := utils.GenerateRandomString()

	now := time.Now()

	tokenMap, err := repository.CreateNewToken(sessionToken, token.AccessToken, &now)

	if err != nil {
		return "", err
	}

	return tokenMap.SessionToken, nil
}

func GetAuthToken(sessionToken string) (string, error) {
	tokenMap, err := repository.GetToken(sessionToken)

	if err != nil {
		return "", err
	}

	return tokenMap.AccessToken, nil
}
