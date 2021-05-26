package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namchee/go-oauth2-client/services"
)

func GetName(ctx *gin.Context) {
	auth, err := ctx.Request.Cookie("Authorization")

	if err != nil || auth.Value == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "No token, no access lmao",
		})

		return
	}

	token, err := services.MapToken(auth.Value)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Cannot authorize request",
		})

		return
	}

	name, err := services.GetUsername(token.AccessToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Cannot get username",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": map[string]string{
			"name": name,
		},
		"error": nil,
	})
}
