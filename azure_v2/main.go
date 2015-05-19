package main

import (
	"log"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/rightscale/go_middleware"

	// load app files
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"
	"github.com/rightscale/self-service-plugins/azure_v2/resources"
)

func main() {
	// Serve
	s := HttpServer()
	log.Printf("Azure plugin - listening on %s\n", *config.ListenFlag)
	s.Run(*config.ListenFlag)
}

// Factory method for application
// Makes it possible to do integration testing.
func HttpServer() *echo.Echo {

	// Setup middleware
	e := echo.New()
	e.Use(middleware.RequestID)                 // Put that first so loggers can log request id
	e.Use(em.Logger)                            // Log to console
	e.Use(middleware.HttpLogger(config.Logger)) // Log to syslog
	e.Use(am.AzureClientInitializer())

	// Setup routes
	resources.SetupSubscriptionRoutes(e)
	resources.SetupInstanceRoutes(e)
	resources.SetupGroupsRoutes(e)
	resources.SetupStorageAccountsRoutes(e)
	resources.SetupProviderRoutes(e)

	return e
}

// Simple wrapper that returns a echo error from a go error
func Error(err error) *echo.HTTPError {
	return &echo.HTTPError{Error: err}
}
