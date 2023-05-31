package jwalk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	HttpClient  *http.Client
	AuthToken   string
	SessionURL  string
	Session     *Session
	SessionLock sync.Mutex
}

func NewClient(httpClient *http.Client, sessionURL, authToken string) *Client {
	client := &Client{
		HttpClient: httpClient,
		AuthToken:  authToken,
		SessionURL: sessionURL,
	}

	req, err := http.NewRequest("GET", client.SessionURL, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+client.AuthToken)

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Println(resp.Body)
	var session Session
	err = json.NewDecoder(resp.Body).Decode(&session)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	client.SessionLock.Lock()
	client.Session = &session
	client.SessionLock.Unlock()

	return client
}
