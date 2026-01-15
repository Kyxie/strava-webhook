package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
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

func getK8sClient() (*kubernetes.Clientset, *rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, config, err
}

func WebhookHandle(c *gin.Context) {
	var event map[string]interface{}
	if err := c.BindJSON(&event); err != nil {
		log.Println("[Webhook] Invalid JSON:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf("[Webhook] Received payload: %+v\n", event)

	if event["object_type"] == "activity" && event["aspect_type"] == "create" {
		log.Printf("[Webhook] New activity received: %v\n", event["object_id"])

		app := os.Getenv("APP")
		if app == "" {
			app = "strava"
		}
		labelSelector := fmt.Sprintf("app=%s", app)

		go func() {
			err := execInPod(labelSelector, []string{"bin/console", "app:strava:import-data"})
			if err != nil {
				log.Println("[K8s] import-data failed:", err)
				return
			}
			err = execInPod(labelSelector, []string{"bin/console", "app:strava:build-files"})
			if err != nil {
				log.Println("[K8s] build-files failed:", err)
				return
			}
			log.Println("[Webhook] Strava update complete")
		}()
	}

	c.Status(http.StatusOK)
}

func execInPod(labelSelector string, command []string) error {
	clientset, config, err := getK8sClient()
	if err != nil {
		return fmt.Errorf("failed to get k8s client: %v", err)
	}

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil || len(pods.Items) == 0 {
		return fmt.Errorf("pod not found in namespace '%s' with label '%s'", namespace, labelSelector)
	}
	targetPod := pods.Items[0]

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(targetPod.Name).
		Namespace(targetPod.Namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("failed to init executor: %v", err)
	}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}

	return nil
}
