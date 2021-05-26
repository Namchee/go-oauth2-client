package controllers

import (
	"fmt"
	"net/http"

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
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "State mismatched!",
		})

		return
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

		return
	}

	// one week
	// handle the cleanup yourself!
	ctx.SetCookie("Authorization", sessionToken, 60*60*24*7, "/", "localhost:8080", true, true)
	ctx.SetSameSite(http.SameSiteLaxMode)

	ctx.Status(http.StatusNoContent)
}

func HandleLogout(ctx *gin.Context) {
	auth, err := ctx.Request.Cookie("Authorization")

	if err != nil || auth.Value == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "Cannot authenticate request",
		})

		return
	}

	if err := services.Logout(auth.Value); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Failed to logout",
		})

		return
	}

	redirectUrl := "http://localhost:8080/static"

	ctx.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}
