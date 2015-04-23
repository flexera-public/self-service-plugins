package main

import (
	"log"
	"log/syslog"
	"os"

	"github.com/labstack/echo"
	"github.com/rightscale/go-digital-ocean/middleware"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	version = "0.0.1"
)

var (
	app        = kingpin.New("gdo", "Digital Ocean RightScale Self-Service plugin.")
	listenFlag = app.Flag("listen", "Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional").Default(":8080").String()
	dumpFlag   = app.Flag("dump", "Dump HTTP requests and responses made to the DO APIs to STDERR").Bool()

	logger *log.Logger // Global syslog logger
)

func main() {
	initLogger()

	// Parse command line
	app.Version(version)
	app.Parse(os.Args[1:])

	// Setup middleware
	e := echo.New()
	e.Use(middleware.RequestID)                      // Put that first so loggers can log request id
	e.Use(echo.Logger)                               // Log to console
	e.Use(middleware.HttpLogger(logger))             // Log to syslog
	e.Use(middleware.DOClientInitializer(*dumpFlag)) // Initialize DigitalOcean API client

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
	e.Run(*listenFlag)
}

// Initialize syslog logger, blow up horribly in case of failure
func initLogger() {
	var err error
	logger, err = syslog.NewLogger(syslog.LOG_NOTICE|syslog.LOG_LOCAL0, 0)
	if err != nil {
		panic("gdo: failed to initialize syslog logger: " + err.Error())
	}
}

// Handle middleware or controller error
func handleError(resp error, c *echo.Context) {
	if logger != nil {
		logger.Printf("ERROR - %s", resp)
	}
	c.String(500, resp.Error())
}
