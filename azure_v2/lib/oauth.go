package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

const (
	tokenEndpoint = "oauth2/token"
)

type authResponse struct {
	Type         string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"` // seconds
	ExpiresOn    string `json:"expires_on"` // seconds
	NotBefore    string `json:"not_before"` // seconds
	Resource     string `json:resource`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:scope`
	Pwd          string `json:pwd_exp`
	PwdUrl       string `json:pwd_url`
}

// Build request to redeem authorization code and get access token
func RequestToken(tenantId string, grantType string, resource string, clientId string, clientSecret string, refreshToken string) (*authResponse, error) {
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", grantType)
	message := grantType
	if resource != "" {
		data.Set("resource", resource)
		message = fmt.Sprintf("%s for resource %s", grantType, resource)
	}
	if refreshToken != "" {
		data.Set("refresh_token", refreshToken)
	}
	path := fmt.Sprintf("%s/%s/%s", config.AuthHost, tenantId, tokenEndpoint)
	fmt.Printf("Requesting %s: %s\n", message, path)
	resp, err := http.PostForm(path, data)
	defer resp.Body.Close()
	if err != nil {
		return nil, GenericException(fmt.Sprintf("Access token refreshing failed: %v", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	if resp.StatusCode >= 400 {
		return nil, GenericException(fmt.Sprintf("Access token refreshing failed: %s", string(body)))
	}
	var response *authResponse

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	return response, nil
}
