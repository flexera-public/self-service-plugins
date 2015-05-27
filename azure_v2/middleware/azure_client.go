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
				accessToken, err = getAccessToken(c)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("URI: %s", c.Request.RequestURI)
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

func getAccessToken(c *echo.Context) (string, error) {
	var grantType string
	var resource string
	// use client specific access token only while app registration
	if c.Request.RequestURI == "/application/register" {
		grantType = "refresh_token"
		resource = ""
	} else {
		grantType = "client_credentials"
		resource = "https://management.core.windows.net/"
	}

	authResponse, err := lib.RequestToken(grantType, resource)
	if err != nil {
		return "", err
	}
	return authResponse.AccessToken, nil
}
