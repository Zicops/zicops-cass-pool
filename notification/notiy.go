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

func SendNotification(title, body, token string) (NotificationOutput, error) {
	var output NotificationOutput
	gqlQuery := fmt.Sprintf(`mutation { sendNotification(title: "%s", body: "%s", token: "%s") { statuscode } }`, title, body, token)
	code, err := PostRequest("http://zicops-notification-server:8080/query", gqlQuery)
	if err != nil {
		return output, err
	}
	output.Statuscode = code
	return output, nil
}

func PostRequest(url, gqlQuery string) (string, error) {
	// make post request to url with body gqlQuery
	// return response, error
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(gqlQuery)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return strconv.Itoa(resp.StatusCode), nil
}
