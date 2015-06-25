package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
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

// Middleware that creates Azure client using credentials in cookie
func AzureClientInitializer() echo.Middleware {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			accessToken, err := getAccessToken(c)
			if err != nil {
				return err
			}

			if c.Request.Header.Get("Content-Type") == "application/json" {
				bodyDecoder := json.NewDecoder(c.Request.Body)
				c.Set("bodyDecoder", bodyDecoder)
			} else {
				return eh.GenericException("Azure plugin supports only \"application/json\" Content-Type.")
			}
			// prepare request params to use from form
			if err := c.Request.ParseForm(); err != nil {
				return eh.GenericException(fmt.Sprintf("Error has occurred while parsing params: %v", err))
			}

			subscriptionId, err := GetCookie(c, "SubscriptionId")
			if err != nil {
				if *config.Env == "development" {
					subscriptionId = *config.SubscriptionIdCred
				}
				if subscriptionId == "" {
					return eh.GenericException("The'SubscriptionId' cookie is required.")
				}
			} else {
				*config.SubscriptionIdCred = subscriptionId
			}

			t := &oauth.Transport{Token: &oauth.Token{AccessToken: accessToken}}
			client := t.Client()
			c.Set("azure", client)
			return h(c)
		}
	}
}

func GetCookie(c *echo.Context, name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", eh.GenericException(fmt.Sprintf("cookie '%s' is missing", cookie))
	}
	return cookie.Value, nil
}

func getAccessToken(c *echo.Context) (string, error) {
	token, err := GetCookie(c, "AccessToken")
	if err != nil {
		return refreshAccessToken(c)
	} else {
		// get access token from cookies
		return token, nil
	}
}

func refreshAccessToken(c *echo.Context) (string, error) {
	creds := new(Credentials)
	var err error
	creds.TenantId, err = GetCookie(c, "TenantId")
	if err != nil {
		if *config.Env == "development" {
			creds.TenantId = *config.TenantIdCred
		}
	}
	creds.ClientId, err = GetCookie(c, "ClientId")
	if err != nil {
		if *config.Env == "development" {
			creds.ClientId = *config.ClientIdCred
		}
	}
	creds.ClientSecret, err = GetCookie(c, "ClientSecret")
	if err != nil {
		if *config.Env == "development" {
			creds.ClientSecret = *config.ClientSecretCred
		}
	}
	creds.RefreshToken, err = GetCookie(c, "RefreshToken")
	if err != nil {
		if *config.Env == "development" {
			creds.RefreshToken = *config.RefreshTokenCred
		}
	}

	if creds.TenantId == "" || creds.ClientId == "" || creds.ClientSecret == "" || creds.RefreshToken == "" {
		return "", eh.GenericException("The credentials are missing in the cookie. Please set 'AccessToken' or combination of 'TenantId', 'ClientId', 'ClientSecret', 'RefreshToken'.")
	}
	// use client specific access token only while app registration
	if c.Request.RequestURI == "/application/register" {
		creds.GrantType = "refresh_token"
		creds.Resource = ""
	} else {
		creds.GrantType = "client_credentials"
		creds.Resource = "https://management.core.windows.net/"
	}
	authResponse, err := creds.RequestToken()
	if err != nil {
		return "", err
	}
	cookie := &http.Cookie{
		Name:  "AccessToken",
		Value: authResponse.AccessToken,
	}
	// set Access Token in the cookie
	http.SetCookie(c.Response.Writer(), cookie)
	return authResponse.AccessToken, nil
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
		return nil, eh.GenericException(fmt.Sprintf("Access token refreshing failed: %v", err))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	if resp.StatusCode >= 400 {
		return nil, eh.GenericException(fmt.Sprintf("Access token refreshing failed: %s", string(body)))
	}
	var response *authResponse

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	return response, nil
}
