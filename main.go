package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"strava-webhook/routes"
)

func main() {
	_ = godotenv.Load()

	port := "8001"

	r := gin.Default()

	r.GET("/webhook", routes.WebhookVerify)
	r.POST("/webhook", routes.WebhookHandle)
	r.GET("/subscription/status", routes.SubscriptionStatus)
	r.POST("/subscription/register", routes.SubscriptionRegister)
	r.POST("/subscription/unregister", routes.SubscriptionUnregister)

	r.Run("0.0.0.0:" + port)
}
