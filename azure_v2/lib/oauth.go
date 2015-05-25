package lib

import (
	"encoding/json"
	"fmt"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	authHost      = "https://login.windows.net"
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
func RequestToken(grantType string, resource string) (*authResponse, error) {
	data := url.Values{}
	data.Set("client_id", *config.ClientIdCred)
	data.Set("client_secret", *config.ClientSecretCred)
	data.Set("refresh_token", *config.RefreshTokenCred)
	data.Set("grant_type", grantType)
	if resource != "" {
		data.Set("resource", resource)
	}
	path := fmt.Sprintf("%s/%s/%s", authHost, *config.TenantIdCred, tokenEndpoint)
	fmt.Printf("Requesting %s: %s\n", grantType, path)
	resp, err := http.PostForm(path, data)
	defer resp.Body.Close()
	if err != nil {
		return nil, GenericException(fmt.Sprintf("failed to build code redeem request: %v", err))
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
