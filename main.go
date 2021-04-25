package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namchee/go-oauth2-client/controllers"
)

func main() {
	router := gin.Default()

	router.LoadHTMLFiles("views/*")
	router.StaticFS("/static", http.Dir("views"))
	router.GET("/oauth/redirect", controllers.HandleOAuthRedirect)

	router.Run(":8080")
}
