package notification

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// package to send gql mutation to http://zicops-notification-server:8080/query

type NotificationOutput struct {
	Statuscode string `json:"statuscode"`
}

func SendNotification(title, body, user_token string, fcm_token string, origin string) (NotificationOutput, error) {
	url := fmt.Sprintf("https://%s/query", origin)
	var output NotificationOutput
	gqlQuery := fmt.Sprintf(`mutation { sendNotification(title: "%s", body: "%s" user_id: ["%s"]) { statuscode } }`, title, body, user_token)
	code, err := PostRequest(url, gqlQuery, user_token, fcm_token)
	if err != nil {
		return output, err
	}
	output.Statuscode = code
	return output, nil
}

func PostRequest(url, gqlQuery string, token string, fcm_token string) (string, error) {
	// make post request to url with body gqlQuery
	// return response, error
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(gqlQuery)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("fcm_token", fcm_token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return strconv.Itoa(resp.StatusCode), nil
}
