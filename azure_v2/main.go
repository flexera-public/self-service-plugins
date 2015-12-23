package main

import (
	"log"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	"github.com/rightscale/go_middleware"

	// load app files
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"
	"github.com/rightscale/self-service-plugins/azure_v2/resources"
)

func main() {
	// Serve
	s := httpServer()
	log.Printf("Azure plugin - listening on %s under %s environment\n", *config.ListenFlag, *config.Env)
	s.Run(*config.ListenFlag)
}

// Factory method for application
// Makes it possible to do integration testing.
func httpServer() *echo.Echo {

	// Setup middleware
	e := echo.New()
	e.Use(middleware.RequestID)                 // Put that first so loggers can log request id
	e.Use(em.Logger())                          // Log to console
	e.Use(middleware.HttpLogger(config.Logger)) // Log to syslog
	e.Use(am.AzureClientInitializer())
	e.Use(em.Recover())

	if config.DebugMode {
		e.SetDebug(true)
	}

	e.SetHTTPErrorHandler(eh.AzureErrorHandler(e)) // override default error handler

	// Setup routes
	prefix := e.Group(*config.AppPrefix) // added prefix to use multiple nginx location on one SS box
	resources.SetupSubscriptionRoutes(prefix)
	resources.SetupInstanceRoutes(prefix)
	resources.SetupGroupsRoutes(prefix)
	resources.SetupStorageAccountsRoutes(prefix)
	resources.SetupProviderRoutes(prefix)
	resources.SetupNetworkRoutes(prefix)
	resources.SetupSubnetsRoutes(prefix)
	resources.SetupIPAddressesRoutes(prefix)
	resources.SetupAuthRoutes(prefix)
	resources.SetupNetworkInterfacesRoutes(prefix)
	resources.SetupImageRoutes(prefix)
	resources.SetupOperationRoutes(prefix)
	resources.SetupAvailabilitySetRoutes(prefix)
	resources.SetupNetworkSecurityGroupRoutes(prefix)
	resources.SetupNetworkSecurityGroupRuleRoutes(prefix)
	resources.SetupInstanceTypesRoutes(prefix)
	resources.SetupRouteTablesRoutes(prefix)

	return e
}
