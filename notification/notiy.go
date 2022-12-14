package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// package to send gql mutation to http://zicops-notification-server:8080/query

type NotificationOutput struct {
	Statuscode string `json:"statuscode"`
}

type results struct {
	MessageId string `json:"message_ids"`
}

type respBody struct {
	Multicast_id  int `json:"multicast_id"`
	Success       int `json:"success"`
	Failure       int `json:"failure"`
	Canonical_ids int `json:"canonical_id"`
	Results       []results
}

func SendNotification(title, body, user_ids string, user_token string, fcm_token string, origin string) (NotificationOutput, error) {
	url := fmt.Sprintf("https://%s/query", origin)
	var output NotificationOutput
	gqlQuery := fmt.Sprintf(`mutation { sendNotification( notification: { title: "%s", body: "%s", user_id: ["%s"] } ) { statuscode } }`, title, body, user_ids)
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
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(gqlQuery)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("fcm-token", fcm_token)
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		//log.Println(err)
		return "", err
	}
	var res respBody
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Println("Error while receiving response from the server")
		log.Println(err)
	}

	defer resp.Body.Close()
	return strconv.Itoa(res.Success), nil
}
