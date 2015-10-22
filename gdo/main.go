package main

import (
	"log"
	"log/syslog"
	"os"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"
	gdm "github.com/rightscale/gdo/middleware"
	"github.com/rightscale/go_middleware"
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
	// Parse command line
	app.Version(version)
	app.Parse(os.Args[1:])

	// Serve
	s := HttpServer()
	log.Printf("gdo - listening on %s\n", *listenFlag)
	s.Run(*listenFlag)
}

// Factory method for application
// Makes it possible to do integration testing.
func HttpServer() *echo.Echo {
	// Initialize global syslog logger
	if l, err := syslog.NewLogger(syslog.LOG_NOTICE|syslog.LOG_LOCAL0, 0); err != nil {
		panic("gdo: failed to initialize syslog logger: " + err.Error())
	} else {
		logger = l
	}

	// Setup middleware
	e := echo.New()
	e.Use(middleware.RequestID)               // Put that first so loggers can log request id
	e.Use(em.Logger)                          // Log to console
	e.Use(middleware.HttpLogger(logger))      // Log to syslog
	e.Use(gdm.DOClientInitializer(*dumpFlag)) // Initialize DigitalOcean API client

	// Setup routes
	SetupDropletsRoutes(e.Group("/droplets"))
	SetupDropletActionsRoutes(e.Group("/droplets/:id/actions"))
	SetupImagesRoutes(e.Group("/images"))
	SetupImageActionsRoutes(e.Group("/images/:id/actions"))
	SetupActionsRoutes(e.Group("/actions"))
	SetupKeysRoutes(e.Group("/keys"))
	SetupRegionsRoutes(e.Group("/regions"))
	SetupSizesRoutes(e.Group("/sizes"))

	// We're done
	return e
}

// Simple wrapper that returns a echo error from a go error
func Error(err error) *echo.HTTPError {
	return &echo.HTTPError{Error: err}
}
