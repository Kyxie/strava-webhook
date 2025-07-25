package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var apiURL = "https://www.strava.com/api/v3/push_subscriptions"

func SubscriptionStatus(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
			"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
		}).
		Get(apiURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode(), "application/json", resp.Body())
}

func SubscriptionRegister(c *gin.Context) {
	client := resty.New()

	checkResp, err := client.R().
		SetQueryParams(map[string]string{
			"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
			"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
		}).
		SetResult([]map[string]interface{}{}).
		Get(apiURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subscriptions := *checkResp.Result().(*[]map[string]interface{})
	if len(subscriptions) > 0 {
		id := subscriptions[0]["id"]
		log.Printf("[Subscription] Existing webhook found (ID: %v), deleting", id)
		_, _ = client.R().
			SetQueryParams(map[string]string{
				"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
				"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
			}).
			Delete(apiURL + "/" + fmt.Sprintf("%v", id))
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
			"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
			"callback_url":  os.Getenv("STRAVA_CALLBACK_URL"),
			"verify_token":  os.Getenv("STRAVA_VERIFY_TOKEN"),
		}).
		Post(apiURL)

	if err != nil {
		log.Println("[Subscription] Registration failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("[Subscription] Webhook registered")
	c.Data(resp.StatusCode(), "application/json", resp.Body())
}

func SubscriptionUnregister(c *gin.Context) {
	client := resty.New()

	checkResp, err := client.R().
		SetQueryParams(map[string]string{
			"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
			"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
		}).
		SetResult([]map[string]interface{}{}).
		Get(apiURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	subscriptions := *checkResp.Result().(*[]map[string]interface{})
	if len(subscriptions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No active webhook subscription"})
		return
	}

	id := subscriptions[0]["id"]
	_, _ = client.R().
		SetQueryParams(map[string]string{
			"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
			"client_secret": os.Getenv("STRAVA_CLIENT_SECRET"),
		}).
		Delete(apiURL + "/" + fmt.Sprintf("%v", id))

	log.Printf("[Subscription] Webhook ID %v unregistered", id)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Webhook ID %v unregistered", id)})
}
