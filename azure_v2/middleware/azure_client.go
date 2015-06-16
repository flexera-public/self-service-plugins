package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

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
				return lib.GenericException("Azure plugin supports only \"application/json\" Content-Type.")
			}
			// prepare request params to use from form
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
	token, err := lib.GetCookie(c, "AccessToken")
	if err != nil {
		return refreshAccessToken(c)
	} else {
		// get access token from cookies
		return token, nil
	}
}

func refreshAccessToken(c *echo.Context) (string, error) {
	tenantId, err := lib.GetCookie(c, "TenantId")
	if err != nil {
		if *config.Env == "development" {
			tenantId = *config.TenantIdCred
		}
	}
	clientId, err := lib.GetCookie(c, "ClientId")
	if err != nil {
		if *config.Env == "development" {
			clientId = *config.ClientIdCred
		}
	}
	clientSecret, err := lib.GetCookie(c, "ClientSecret")
	if err != nil {
		if *config.Env == "development" {
			clientSecret = *config.ClientSecretCred
		}
	}
	refreshToken, err := lib.GetCookie(c, "RefreshToken")
	if err != nil {
		if *config.Env == "development" {
			refreshToken = *config.RefreshTokenCred
		}
	}

	if tenantId == "" || clientId == "" || clientSecret == "" || refreshToken == "" {
		return "", lib.GenericException("The credentials are missing in the cookie. Please set 'AccessToken' or combination of 'TenantId', 'ClientId', 'ClientSecret', 'RefreshToken'.")
	}
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
	authResponse, err := lib.RequestToken(tenantId, grantType, resource, clientId, clientSecret, *config.RefreshTokenCred)
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
