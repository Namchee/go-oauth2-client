package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namchee/go-oauth2-client/controllers"
)

func main() {
	router := gin.Default()

	router.StaticFS("static", http.Dir("static"))

	router.GET("/login", controllers.HandleLogin)
	router.GET("/oauth/redirect", controllers.HandleOAuthCallback)

	router.POST("/logout", controllers.HandleLogout)

	router.GET("/api/name", controllers.GetName)

	router.Run(":8080")
}
