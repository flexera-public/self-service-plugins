#Azure plugin

Azure plugin is a Go application which serves REST HTTP requests and provides ability to iterate with Azure V2 API

##Build application

```
go build -o azure_plugin
```

## Usage

```
azure_plugin --help
usage: azure_plugin [<flags>] <client> <secret> <resource> <subscription> <refresh_token>

Azure V2 RightScale Self-Service plugin.

Flags:
  --help            Show help.
  --listen=":8080"  Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional
  --version         Show application version.

Args:
  <client>         The client id of the application that is registered in Azure Active Directory.
  <secret>         The client key of the application that is registered in Azure Active Directory.
  <resource>       The App ID URI of the web API (secured resource).
  <subscription>   The client subscription id.
  <refresh_token>  The token used for refreshing access token.
```

##Run tests

```
go test ./resources
```