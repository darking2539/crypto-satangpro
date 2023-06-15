package utils

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func LineNotify(message string) {

	accessToken := os.Getenv("LINE_NOTIFY_TOKEN")

	if accessToken == "" {
		return
	}

	formData := url.Values{}
	formData.Set("message", message)

	payload := strings.NewReader(formData.Encode())

	// Create a request to the Line Messaging API
	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", payload)
	if err != nil {
		log.Println(err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send Line notification. Status:", resp.Status)
		return
	}

}
