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

			subscriptionId, err := lib.GetCookie(c, "SubscriptionId")
			if err != nil {
				if *config.Env == "development" {
					subscriptionId = *config.SubscriptionIdCred
				}
				if subscriptionId == "" {
					return lib.GenericException("The'SubscriptionId' cookie is required.")
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
	creds := new(lib.Credentials)
	var err error
	creds.TenantId, err = lib.GetCookie(c, "TenantId")
	if err != nil {
		if *config.Env == "development" {
			creds.TenantId = *config.TenantIdCred
		}
	}
	creds.ClientId, err = lib.GetCookie(c, "ClientId")
	if err != nil {
		if *config.Env == "development" {
			creds.ClientId = *config.ClientIdCred
		}
	}
	creds.ClientSecret, err = lib.GetCookie(c, "ClientSecret")
	if err != nil {
		if *config.Env == "development" {
			creds.ClientSecret = *config.ClientSecretCred
		}
	}
	creds.RefreshToken, err = lib.GetCookie(c, "RefreshToken")
	if err != nil {
		if *config.Env == "development" {
			creds.RefreshToken = *config.RefreshTokenCred
		}
	}

	if creds.TenantId == "" || creds.ClientId == "" || creds.ClientSecret == "" || creds.RefreshToken == "" {
		return "", lib.GenericException("The credentials are missing in the cookie. Please set 'AccessToken' or combination of 'TenantId', 'ClientId', 'ClientSecret', 'RefreshToken'.")
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
