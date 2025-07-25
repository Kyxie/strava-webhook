package routes

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func WebhookVerify(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == os.Getenv("STRAVA_VERIFY_TOKEN") {
		log.Println("[Webhook] Verified successfully")
		c.JSON(http.StatusOK, gin.H{"hub.challenge": challenge})
	} else {
		log.Println("[Webhook] Verification failed")
		c.String(http.StatusForbidden, "Forbidden")
	}
}

func WebhookHandle(c *gin.Context) {
	var event map[string]interface{}
	if err := c.BindJSON(&event); err != nil {
		log.Println("[Webhook] Invalid JSON:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if event["object_type"] == "activity" && event["event_type"] == "create" {
		log.Printf("[Webhook] New activity received: %v\n", event["object_id"])

		go func() {
			if err := run("docker", "exec", "strava", "bin/console", "app:strava:import-data"); err != nil {
				log.Println("[Webhook] import-data failed:", err)
				return
			}
			if err := run("docker", "exec", "strava", "bin/console", "app:strava:build-files"); err != nil {
				log.Println("[Webhook] build-files failed:", err)
				return
			}
			log.Println("[Webhook] Strava update complete")
		}()
	}

	c.Status(http.StatusOK)
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Command] %s %v failed\nOutput: %s\n", name, args, string(output))
		return err
	}
	log.Printf("[Command] %s %v succeeded\nOutput: %s\n", name, args, string(output))
	return nil
}
