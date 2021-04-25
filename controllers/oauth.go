package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	Utils "github.com/namchee/go-oauth2-client/utils"
)

type OAuth2Response struct {
	AccessToken string `json:"access_token"`
}

func HandleOAuthRedirect(ctx *gin.Context) {
	requestCode := ctx.Query("code")

	if requestCode == "" {
		ctx.JSON(
			http.StatusForbidden, gin.H{
				"error": "Request token is empty!",
			},
		)
	}

	credentials := Utils.GetCredentials()

	reqUrl := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		credentials.ClientId,
		credentials.ClientSecret,
		requestCode,
	)

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create HTTP request",
		})
	}

	req.Header.Set("Accept", "application/json")

	res, err := httpClient.Do(req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send HTTP request",
		})
	}

	defer res.Body.Close()

	var jsonResponse OAuth2Response

	if err := json.NewDecoder(res.Body).Decode(&jsonResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse JSON response",
		})
	}

	authorizedUrl := fmt.Sprintf("http://localhost:8080/static/authorized.html?access_token=%s", jsonResponse.AccessToken)

	ctx.Redirect(http.StatusFound, authorizedUrl)
}
