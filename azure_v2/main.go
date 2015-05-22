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
	log.Printf("Azure plugin - listening on %s under %s environment\n", *config.ListenFlag, *config.Env)
	s.Run(*config.ListenFlag)
}

// Factory method for application
// Makes it possible to do integration testing.
func HttpServer() *echo.Echo {

	// Setup middleware
	e := echo.New()
	e.Use(middleware.RequestID)                 // Put that first so loggers can log request id
	e.Use(em.Logger())                          // Log to console
	e.Use(middleware.HttpLogger(config.Logger)) // Log to syslog
	e.Use(am.AzureClientInitializer())

	if config.DebugMode {
		e.SetDebug(true)
	}

	e.SetHTTPErrorHandler(AzureErrorHandler(e)) // override default error handler

	// Setup routes
	resources.SetupSubscriptionRoutes(e)
	resources.SetupInstanceRoutes(e)
	resources.SetupGroupsRoutes(e)
	resources.SetupStorageAccountsRoutes(e)
	resources.SetupProviderRoutes(e)
	resources.SetupNetworkRoutes(e)
	resources.SetupSubnetsRoutes(e)
	resources.SetupIpAddressesRoutes(e)

	return e
}
