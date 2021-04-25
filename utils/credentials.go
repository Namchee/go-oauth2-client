package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Credentials struct {
	ClientId     string
	ClientSecret string
}

func (c *Credentials) isNil() bool {
	return c.ClientId == "" && c.ClientSecret == ""
}

var (
	credentials Credentials
)

func GetCredentials() Credentials {
	if credentials.isNil() {
		godotenv.Load()

		clientId, isSet := os.LookupEnv("CLIENT_ID")

		if !isSet {
			log.Fatalln("Client ID doesn't exist")
		}

		clientSecret, isSet := os.LookupEnv("CLIENT_SECRET")

		if !isSet {
			log.Fatalln("Client secret doesn't exist")
		}

		credentials = Credentials{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		}
	}

	return credentials
}
