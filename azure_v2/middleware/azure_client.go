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
)

const (
	authHost       = "https://login.windows.net"
	tokenEndpoint  = "/oauth2/token"
	CredCookieName = "AccessToken"
)

var (
	accessToken  string
	authResponse struct {
		Type         string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"` // seconds
		ExpiresOn    int64  `json:"expires_on"` // seconds
		NotBefore    int64  `json:"not_before"` // seconds
		Resource     string `json:resource`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:scope`
		Pwd          int64  `json:pwd_exp`
		PwdUrl       string `json:pwd_url`
	}
)

// Middleware that creates Azure client using credentials in cookie
func AzureClientInitializer() echo.Middleware {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) *echo.HTTPError {
			token, err := c.Request.Cookie(CredCookieName)
			if err != nil {
				resp, err := refreshAccessToken()
				if err != nil {
					return &echo.HTTPError{
						Error: fmt.Errorf("failed to build code redeem request: %v", err),
						Code:  400,
					}
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Errorf("failed to load response body: %s", err)
				}
				if resp.StatusCode >= 400 {
					return &echo.HTTPError{
						Error: fmt.Errorf("Access token refreshing failed: %s", string(body)),
						Code:  resp.StatusCode,
					}
				}

				if err = json.Unmarshal(body, &authResponse); err != nil {
					fmt.Errorf("got bad response from server: %q", body)
				}
				accessToken = authResponse.AccessToken
			} else {
				// get access token from cookies
				accessToken = token.Value
			}

			// prepare request params to use
			if err := c.Request.ParseForm(); err != nil {
				return &echo.HTTPError{
					Error: fmt.Errorf("parseForm(): %v", err),
					Code:  400,
				}
			}

			t := &oauth.Transport{Token: &oauth.Token{AccessToken: accessToken}}
			client := t.Client()
			c.Set("azure", client)
			return h(c)
		}
	}
}

// Retrieve client initialized by middleware, send error response if not found
// This function should be used by controller actions that need to use the client
func GetAzureClient(c *echo.Context) (*http.Client, *echo.HTTPError) {
	client, _ := c.Get("azure").(*http.Client)
	if client == nil {
		return nil, &echo.HTTPError{Error: fmt.Errorf("failed to retrieve Azure client, check middleware")}
	}
	return client, nil
}

func GetCookie(c *echo.Context, name string) (*http.Cookie, *echo.HTTPError) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return nil, &echo.HTTPError{
			Error: fmt.Errorf("cookie '%s' is missing", cookie),
			Code:  400,
		}
	}
	return cookie, nil
}

// Build request to redeem authorization code and get access token
func refreshAccessToken() (*http.Response, error) {
	data := url.Values{}
	data.Set("client_id", *config.ClientIdCred)
	data.Set("client_secret", *config.ClientSecretCred)
	data.Set("refresh_token", *config.RefreshTokenCred)
	data.Set("grant_type", "refresh_token")
	fmt.Printf("Refreshing access token ...\n")
	resp, err := http.PostForm(authHost+"/common"+tokenEndpoint, data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
