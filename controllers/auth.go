package controllers

import (
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

func HandleLogin(ctx *gin.Context) {
	url := services.GetLoginUrl()

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleOAuthCallback(ctx *gin.Context) {
	var callbackValue OAuth2CallbackForm

	callbackValue.State = ctx.Query("state")
	callbackValue.Code = ctx.Query("code")

	if callbackValue.State == "" || !services.ValidateState(callbackValue.State) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "State mismatched",
		})

		return
	}

	if callbackValue.Code == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "Authorization code is required",
		})
	}

	sessionToken, err := services.GetSessionToken(callbackValue.Code)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Failed to create session token",
		})

		return
	}

	// one week
	// handle the cleanup yourself!
	// secure is set to false for localhost debugging
	ctx.SetCookie("Authorization", sessionToken, 60*60*24*7, "/", "localhost:8080", false, true)
	ctx.SetSameSite(http.SameSiteLaxMode)

	ctx.Redirect(http.StatusTemporaryRedirect, "/static/authorized")
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

	ctx.Redirect(http.StatusTemporaryRedirect, "/static")
}
