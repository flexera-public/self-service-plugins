package config

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
	"log/syslog"
	"os"
)

const (
	version = "0.0.1"
	// APIVersion is a default Azure API version
	APIVersion = "2014-12-01-Preview"
	// MediaType is default media type for requests to the Azure cloud
	MediaType = "application/json"
	// UserAgent is a RS request sign
	UserAgent = "RightScale Self-Service Plugin"
)

var (
	app = kingpin.New("azure_plugin", "Azure V2 RightScale Self-Service plugin.")
	// ListenFlag is a hostname and port to listen
	ListenFlag = app.Flag("listen", "Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional.").Default("localhost:8080").String()
	// Env is environment name
	Env = app.Flag("env", "Environment name: 'development' (default) or 'production'.").Default("development").String()
	// AppPrefix is URL prefix
	AppPrefix = app.Flag("prefix", "URL prefix.").Default("/azure_plugin").String()
	// ClientIDCred is the client id of the application that is registered in Azure Active Directory.
	ClientIDCred = app.Arg("client", "The client id of the application that is registered in Azure Active Directory.").String()
	// ClientSecretCred is the client key of the application that is registered in Azure Active Directory.
	ClientSecretCred = app.Arg("secret", "The client key of the application that is registered in Azure Active Directory.").String()
	// SubscriptionIDCred is the client subscription id.
	SubscriptionIDCred = app.Arg("subscription", "The client subscription id.").String()
	// TenantIDCred is Azure Active Directory indentificator.
	TenantIDCred = app.Arg("tenant", "Azure Active Directory indentificator.").String()
	// RefreshTokenCred is the token used for refreshing access token.
	RefreshTokenCred = app.Arg("refresh_token", "The token used for refreshing access token.").String()
	// BaseURL is Azure cloud endpoint...set base url as variable to be able to modify it in the specs
	BaseURL = "https://management.azure.com"
	// GraphURL is the endpoint to Graph Azure service
	GraphURL = "https://graph.windows.net"
	// AuthHost is endpoint to authentication Azure service
	AuthHost = "https://login.windows.net"
	// Logger is Global syslog logger
	Logger *log.Logger
	// DebugMode is used to manage debug mode
	DebugMode = false
)

func init() {
	// Parse command line
	app.Version(version)
	app.Parse(os.Args[1:])

	// Initialize global syslog logger
	if l, err := syslog.NewLogger(syslog.LOG_NOTICE|syslog.LOG_LOCAL0, 0); err != nil {
		panic("azure: failed to initialize syslog logger: " + err.Error())
	} else {
		Logger = l
	}

	switch *Env {
	case "development":
		// add development specific settings here
		DebugMode = true
	case "production":
		// add production specific settings here
		// example: *ListenFlag = "rightscale.com:80"
	default:
		panic("Unknown environmental name: " + *Env)
	}

}
