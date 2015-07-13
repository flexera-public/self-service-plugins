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

const tokenEndpoint = "oauth2/token"

// Credentials represents set of creds required for Azure authentication
type Credentials struct {
	TenantID     string `json:"tenant"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Subscription string `json:"subscription"`
	GrantType    string
	Resource     string
	RefreshToken string
}

// AuthResponse represents creds gotten from cloud
type AuthResponse struct {
	Type         string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"` // seconds
	ExpiresOn    string `json:"expires_on"` // seconds
	NotBefore    string `json:"not_before"` // seconds
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Pwd          string `json:"pwd_exp"`
	PwdURL       string `json:"pwd_url"`
}

// AzureClientInitializer is a middleware that creates Azure client and handles credentials
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

			subscriptionID, err := getCookie(c, "SubscriptionID")
			if err != nil {
				if *config.Env == "development" {
					subscriptionID = *config.SubscriptionIDCred
				}
				if subscriptionID == "" {
					return eh.GenericException("The'SubscriptionID' cookie is required.")
				}
			} else {
				*config.SubscriptionIDCred = subscriptionID
			}

			t := &oauth.Transport{Token: &oauth.Token{AccessToken: accessToken}}
			client := t.Client()
			c.Set("azure", client)
			return h(c)
		}
	}
}

func getCookie(c *echo.Context, name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", eh.GenericException(fmt.Sprintf("cookie '%s' is missing", cookie))
	}
	return cookie.Value, nil
}

func getAccessToken(c *echo.Context) (string, error) {
	token, err := getCookie(c, "AccessToken")
	if err != nil {
		return refreshAccessToken(c)
	}
	// get access token from cookies
	return token, nil
}

func refreshAccessToken(c *echo.Context) (string, error) {
	creds := new(Credentials)
	var err error
	creds.TenantID, err = getCookie(c, "TenantID")
	if err != nil {
		if *config.Env == "development" {
			creds.TenantID = *config.TenantIDCred
		}
	}
	*config.TenantIDCred = creds.TenantID
	creds.ClientID, err = getCookie(c, "ClientID")
	if err != nil {
		if *config.Env == "development" {
			creds.ClientID = *config.ClientIDCred
		}
	}
	*config.ClientIDCred = creds.ClientID
	creds.ClientSecret, err = getCookie(c, "ClientSecret")
	if err != nil {
		if *config.Env == "development" {
			creds.ClientSecret = *config.ClientSecretCred
		}
	}
	*config.ClientSecretCred = creds.ClientSecret
	creds.RefreshToken, err = getCookie(c, "RefreshToken")
	if err != nil {
		if *config.Env == "development" {
			creds.RefreshToken = *config.RefreshTokenCred
		}
	}
	*config.RefreshTokenCred = creds.RefreshToken

	if creds.TenantID == "" || creds.ClientID == "" || creds.ClientSecret == "" || creds.RefreshToken == "" {
		return "", eh.GenericException("The credentials are missing in the cookie. Please set 'AccessToken' or combination of 'TenantID', 'ClientID', 'ClientSecret', 'RefreshToken'.")
	}
	// use client specific access token only while app registration
	//TODO: use regexp here
	if c.Request.RequestURI == *config.AppPrefix+"/application/register" || c.Request.RequestURI == *config.AppPrefix+"/application/unregister" {
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
	// set Access Token in the cookie
	http.SetCookie(c.Response.Writer(), &http.Cookie{
		Name:  "AccessToken",
		Value: authResponse.AccessToken,
	})
	// set ExpiresOn in the cookie
	http.SetCookie(c.Response.Writer(), &http.Cookie{
		Name:  "ExpiresOn",
		Value: authResponse.ExpiresOn,
	})
	return authResponse.AccessToken, nil
}

// RequestToken builds request to redeem authorization code and get access token
func (c *Credentials) RequestToken() (*AuthResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
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
	path := fmt.Sprintf("%s/%s/%s", config.AuthHost, c.TenantID, tokenEndpoint)
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
	var response *AuthResponse

	if err = json.Unmarshal(body, &response); err != nil {
		return nil, eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	return response, nil
}
