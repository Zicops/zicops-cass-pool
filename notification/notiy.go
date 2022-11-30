package notification

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// package to send gql mutation to http://zicops-notification-server:8080/query

type NotificationOutput struct {
	Statuscode string `json:"statuscode"`
}

func SendNotification(title, body, token, fcmToken string) (NotificationOutput, error) {
	var output NotificationOutput
	// url from env
	url := os.Getenv("NOTIFICATION_URL")
	if url == "" {
		url = "https://demo.zicops.com/ns/query"
	}
	gqlQuery := fmt.Sprintf(`mutation { sendNotification(notification: { title: "%s", body: "%s" } ) { statuscode } }`, title, body)
	code, err := PostRequest(url, token, fcmToken, gqlQuery)
	if err != nil {
		return output, err
	}
	output.Statuscode = code
	return output, nil
}

func PostRequest(url, token, fcmToken, gqlQuery string) (string, error) {
	// make post request to url with body gqlQuery
	// return response, error
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(gqlQuery)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("fcm-token", fcmToken)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return strconv.Itoa(resp.StatusCode), nil
}
