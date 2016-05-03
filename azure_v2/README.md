#Azure plugin

Azure plugin is a Go application which serves REST HTTP requests and provides ability to iterate with Azure V2 API
TODO: investigate new azure repo - https://github.com/Azure/azure-sdk-for-go/tree/master/arm

##Build application

```
go build -o azure_plugin
```

## Usage

```
azure_plugin --help
usage: azure_plugin [<flags>] [<client> [<secret> [<subscription> [<tenant> [<refresh_token>]]]]]

Azure V2 RightScale Self-Service plugin.

Flags:
  --help               Show help.
  --listen="localhost:8080"
                       Hostname and port to listen on, e.g. 'localhost:8080' - hostname is optional.
  --env="development"  Environment name: 'development' (default) or 'production'.
  --prefix="/azure_plugin"
                       URL prefix.
  --version            Show application version.

Args:
  [<client>]         The client id of the application that is registered in Azure Active Directory.
  [<secret>]         The client key of the application that is registered in Azure Active Directory.
  [<subscription>]   The client subscription id.
  [<tenant>]         Azure Active Directory indentificator.
  [<refresh_token>]  The token used for refreshing access token.
```

##New cloud registration
First step of cloud registration is registering RS application in the client Active Directory
in order to get ability to use application specific access token.
curl -v 'http://localhost:8080/application/register'
Note: don't forget about creds in the cookies (see below)

##Unregister application - "Disconnect" Subscription from Application
From Azure docs: "Just as you enabled users to connect their subscriptions to your application, you must allow then to disconnect subscriptions too. From an access management point of view, disconnect means removing the role assignment that the applications service principal has on the subscription."
curl -v -b ... 'http://localhost:8080/application/unregister'

##Make requests
With no access token passed in the cookies
curl -v -b "TenantID=...;ClientID=...;ClientSecret=...;SubscriptionID=...;RefreshToken=..." 'http://localhost:8080/instances'
This kind of call will go through azure oauth to get application specific access token.
Note: application should be registered in advance

Pass access token in the cookies
curl -v -b "AccessToken=eyJ0eXAiOiJKV1QiLCJhbGci...;SubscriptionID=..." 'http://localhost:8080/instances'
Note: could be used either user or app specific access token but take into account that plugin doesn't refresh token automatically

##Run tests

```
go test ./resources
```

## Nginx config for SS box
```
location /azure_v2 {
    proxy_connect_timeout       600;
    proxy_send_timeout          600;
    proxy_read_timeout          600;
    proxy_pass http://azure_v2;
}

upstream azure_v2  {
    server localhost:8083;
}
```

## To run on SS box

```
tar zxvf binary/azure_v2-linux-amd64.tgz
nohup ./azure_v2 --listen="localhost:8083" --prefix="/azure_v2" > ./azure_v2.log 2>&1 &
```