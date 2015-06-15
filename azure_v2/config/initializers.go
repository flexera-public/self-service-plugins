package config

import (
	"gopkg.in/alecthomas/kingpin.v1"
	"log"
	"log/syslog"
	"os"
)

const (
	version    = "0.0.1"
	ApiVersion = "2014-12-01-Preview"
	MediaType  = "application/json"
	UserAgent  = "RightScale Self-Service Plugin"
)

var (
	app                = kingpin.New("azure_plugin", "Azure V2 RightScale Self-Service plugin.")
	ListenFlag         = app.Flag("listen", "Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional.").Default("localhost:8080").String()
	Env                = app.Flag("env", "Environment name: 'development' (default) or 'production'.").Default("development").String()
	ClientIdCred       = app.Arg("client", "The client id of the application that is registered in Azure Active Directory.").Required().String()
	ClientSecretCred   = app.Arg("secret", "The client key of the application that is registered in Azure Active Directory.").Required().String()
	ResourceCred       = app.Arg("resource", "The App ID URI of the web API (secured resource).").Required().String()
	SubscriptionIdCred = app.Arg("subscription", "The client subscription id.").Required().String()
	TenantIdCred       = app.Arg("tenant", "Azure Active Directory indentificator.").Required().String()
	RefreshTokenCred   = app.Arg("refresh_token", "The token used for refreshing access token.").Required().String()
	// set base url as variable to be able to modify it in the specs
	BaseUrl   = "https://management.azure.com"
	AuthHost  = "https://login.windows.net"
	GraphUrl  = "https://graph.windows.net"
	Logger    *log.Logger // Global syslog logger
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
