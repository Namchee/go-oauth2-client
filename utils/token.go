package utils

import "github.com/google/uuid"

func GenerateRandomString() (string, error) {
	sessionToken, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return sessionToken.String(), nil
}
