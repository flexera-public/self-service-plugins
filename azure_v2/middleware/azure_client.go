package middleware

import (
	"fmt"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	CredCookieName = "AccessToken"
)

var (
	accessToken string
)

// Middleware that creates Azure client using credentials in cookie
func AzureClientInitializer() echo.Middleware {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, err := c.Request.Cookie(CredCookieName)
			if err != nil {
				authResponse, err := lib.RequestToken("refresh_token", "")
				if err != nil {
					return err
				}
				accessToken = authResponse.AccessToken
			} else {
				// get access token from cookies
				accessToken = token.Value
			}

			// prepare request params to use
			if err := c.Request.ParseForm(); err != nil {
				return lib.GenericException(fmt.Sprintf("Error has occurred while parsing params: %v", err))
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
func GetAzureClient(c *echo.Context) (*http.Client, error) {
	client, _ := c.Get("azure").(*http.Client)
	if client == nil {
		return nil, lib.GenericException(fmt.Sprintf("failed to retrieve Azure client, check middleware"))
	}
	return client, nil
}

func GetCookie(c *echo.Context, name string) (*http.Cookie, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("cookie '%s' is missing", cookie))
	}
	return cookie, nil
}
