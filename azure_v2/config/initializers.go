package config

import (
	"os"
	"log"
	"log/syslog"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	version = "0.0.1"
	ApiVersion = "2014-12-01-Preview"
	BaseUrl = "https://management.azure.com"
)

var (
	app        = kingpin.New("azure", "Azure V2 RightScale Self-Service plugin.")
	ListenFlag = app.Flag("listen", "Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional").Default(":8080").String()

	Logger *log.Logger // Global syslog logger
)

func init(){
	// Parse command line
	app.Version(version)
	app.Parse(os.Args[1:])

	// Initialize global syslog logger
	if l, err := syslog.NewLogger(syslog.LOG_NOTICE|syslog.LOG_LOCAL0, 0); err != nil {
		panic("azure: failed to initialize syslog logger: " + err.Error())
	} else {
		Logger = l
	}
}