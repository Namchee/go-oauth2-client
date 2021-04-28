package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/namchee/go-oauth2-client/services"
)

type OAuth2CallbackForm struct {
	State string `json:"state"`
	Code  string `json:"code"`
}

type OAuth2AccessResponse struct {
	AccessToken string `json:"access_token"`
}

type TokenRequestForm struct {
	AuthCode string `json:"auth_code"`
}

func HandleLogin(ctx *gin.Context) {
	url := services.GetLoginUrl()

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleOAuthCallback(ctx *gin.Context) {
	var callbackValue OAuth2CallbackForm

	callbackValue.State = ctx.Query("state")
	callbackValue.Code = ctx.Query("code")

	if !services.ValidateState(callbackValue.State) {
		ctx.HTML(http.StatusForbidden, "/unauthorized", gin.H{
			"reason": "State mismatched!",
		})
	}

	redirectUrl := fmt.Sprintf("/static/callback.html?auth_code=%s", callbackValue.Code)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func HandleTokenRequest(ctx *gin.Context) {
	var requestParam TokenRequestForm

	ctx.BindJSON(&requestParam)

	sessionToken, err := services.GetSessionToken(requestParam.AuthCode)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Failed to create session token",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]string{
			"session_token": sessionToken,
		},
		"error": nil,
	})
}

func GetName(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "No token, no access lmao",
		})
	}

	fmt.Println(authHeader)

	tokens := strings.Split(authHeader, " ")
	sessionToken := tokens[1]

	accessToken, err := services.GetAccessToken(sessionToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Cannot authorize request",
		})
	}

	name, err := services.GetUsername(accessToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Cannot get username",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]string{
			"name": name,
		},
		"error": nil,
	})
}
