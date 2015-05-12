package middleware

import (
	"fmt"
	"net/http"
	
	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
)

// Name of cookie created by SS that contains the credentials needed to send API requests to Azure
const (
	credCookieName = "AccessToken"
	SubscriptionCookieName = "SubscriptionId"
)

// Middleware that creates Azure client using credentials in cookie
func AzureClientInitializer() echo.Middleware {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) *echo.HTTPError {
			token, _ := GetCookie(c, credCookieName)
			
			err := c.Request.ParseForm()
			if err != nil {
				fmt.Errorf("parseForm(): %v", err)
			}

			t := &oauth.Transport{Token: &oauth.Token{AccessToken: token.Value}}
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
