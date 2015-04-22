package main

import "github.com/labstack/echo"

func main() {
	e := echo.New()

	// Setup middleware
	e.Use(echo.Logger) // Log to console
	e.Use(HttpLogger)  // Log to syslog
	e.Use(DOClient())  // Initialize DigitalOcean API client

	// Setup error handler
	e.HTTPErrorHandler(handleError)

	// Setup routes
	SetupDropletsRoutes(e.Group("/droplets"))
	SetupDropletActionsRoutes(e.Group("/droplets/:id/actions"))
	SetupImagesRoutes(e.Group("/images"))
	SetupImageActionsRoutes(e.Group("/images/:id/actions"))
	SetupActionsRoutes(e.Group("/actions"))
	SetupKeysRoutes(e.Group("/keys"))
	SetupRegionsRoutes(e.Group("/regions"))
	SetupSizesRoutes(e.Group("/sizes"))

	// Serve
	e.Run(":8080")
}

// Handle middleware or controller error
func handleError(resp error, c *echo.Context) {
	logger.Printf("ERROR - %s", resp)
	c.String(500, "")
}
