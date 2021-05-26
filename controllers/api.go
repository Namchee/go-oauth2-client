package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/namchee/go-oauth2-client/services"
)

func GetName(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")

	if authHeader == "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"data":  nil,
			"error": "No token, no access lmao",
		})
	}

	tokens := strings.Split(authHeader, " ")
	sessionToken := tokens[1]

	token, err := services.MapToken(sessionToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": "Cannot authorize request",
		})
	}

	name, err := services.GetUsername(token.AccessToken)

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
