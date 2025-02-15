package weatherservice

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const mateoMaticsAuthenticationURL = "https://login.meteomatics.com/api/v1/token"

type MateoMaticAPI struct {
	username string
	password string
	token    string

	client *http.Client
}

func NewMateoMaticAPI(username, password string) *MateoMaticAPI {
	m := &MateoMaticAPI{
		username: username,
		password: password,
		client:   http.DefaultClient,
	}

	if m.getOAuthToken() != nil {
		return nil
	}

	return m
}

func (api *MateoMaticAPI) getOAuthToken() error {
	if api.client == nil {
		return fmt.Errorf("client not initialized")
	}

	// Encode username and password in Base64 for Basic Auth
	auth := api.username + ":" + api.password
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))

	// Create HTTP request
	req, err := http.NewRequest("GET", mateoMaticsAuthenticationURL, nil)
	if err != nil {
		return err
	}

	// Add Authorization header
	req.Header.Set("Authorization", "Basic "+authEncoded)

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Parse JSON token response
	var result struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	api.token = result.AccessToken

	return nil

}

func (api *MateoMaticAPI) FetchReportData(ctx context.Context, city string) (*ReportData[Waves], error) {
	return nil, nil
}
