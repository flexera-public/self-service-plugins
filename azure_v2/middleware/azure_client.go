package middleware

import (
	"fmt"

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
