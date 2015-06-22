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

type Credentials struct {
	TenantId     string `json:"tenant"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Subscription string `json:"subscription"`
	GrantType    string
	Resource     string
	RefreshToken string
}

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
func (c *Credentials) RequestToken() (*authResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientId)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", c.GrantType)
	message := c.GrantType
	if c.Resource != "" {
		data.Set("resource", c.Resource)
		message = fmt.Sprintf("%s for resource %s", c.GrantType, c.Resource)
	}
	if c.RefreshToken != "" {
		data.Set("refresh_token", c.RefreshToken)
	}
	path := fmt.Sprintf("%s/%s/%s", config.AuthHost, c.TenantId, tokenEndpoint)
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
