package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/godo"
)

// Global syslog logger
var logger *log.Logger

// Simple HTTP logger middleware
func HttpLogger() echo.MiddlewareFunc {
	var err error
	logger, err = syslog.NewLogger(syslog.LOG_INFO|syslog.LOG_LOCAL0, 0)
	if err != nil {
		panic(err.Error())
	}
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			msg := fmt.Sprintf(`Processing GET "%s" (for %s)`)
			if reqId := c.Request.Header.Get("X-Request-Id"); reqId != "" {
				msg += fmt.Sprintf("Request ID: %s", reqId)
			}
			logger.Print(msg)
			start := time.Now()
			err := h(c)
			if err != nil {
				return err
			}
			elapsed := time.Since(start)
			var status string
			var size int
			if resp := c.Response; resp != nil {
				status = http.StatusText(resp.Status())
				size = resp.Size()
			}
			logger.Printf(`Completed in %s | %s | %d bytes`, elapsed, status, size)
			return nil
		}
	}
}

// Name of cookie created by SS that contains the credentials needed to send API requests to DO
const CredCookieName = "ServiceCred"

// Middleware that creates DO client using credentials in cookie
func DOClient() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			token, err := c.Request.Cookie(CredCookieName)
			if err != nil {
				return err
			}
			t := &oauth.Transport{Token: &oauth.Token{AccessToken: token.Value}}
			client := godo.NewClient(t.Client())
			c.Set("doC", client)
			return h(c)
		}
	}

}

// Retrieve client initialized by middleware, send error response if not found
// This function should be used by controller actions that need to use the client
func GetDOClient(c *echo.Context) (*godo.Client, error) {
	client, _ := c.Get("doC").(*godo.Client)
	if client == nil {
		return nil, fmt.Errorf("failed to retrieve Digital Ocean client, check middleware")
	}
	return client, nil
}
