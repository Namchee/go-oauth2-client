package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namchee/go-oauth2-client/controllers"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("views/*")
	router.StaticFS("static", http.Dir("static"))

	router.GET("/login", controllers.HandleLogin)
	router.POST("/auth/callback", controllers.HandleOAuthCallback)
	router.POST("/auth/token", controllers.HandleTokenRequest)
	router.GET("/api/name")

	router.Run(":8080")
}
