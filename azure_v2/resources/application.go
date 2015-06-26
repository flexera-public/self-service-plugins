package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

const (
	authPath          = "providers/Microsoft.Authorization/roleDefinitions"
	roleContributorID = "b24988ac-6180-42a0-ab88-20f7382dd24c"
)

type servicePrincipal struct {
	// Add more fields if needed
	ObjectID string `json:"objectId"`
}

// SetupAuthRoutes declares routes for Application resource
func SetupAuthRoutes(e *echo.Echo) {
	e.Post("/application/register", registerApp)
}

// Get App-Only Access Token for Azure AD Graph API
func registerApp(c *echo.Context) error {
	creds := am.Credentials{GrantType: "client_credentials", Resource: "https://graph.windows.net/"}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&creds)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	authResponse, err := creds.RequestToken()
	if err != nil {
		return err
	}
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: authResponse.AccessToken}}
	graphClient := t.Client()
	principalID, err := getServicePrincipal(graphClient, &creds)
	if err != nil {
		return err
	}
	return assignRoleToApp(c, principalID, creds.Subscription)
}

func getServicePrincipal(client *http.Client, creds *am.Credentials) (string, error) {
	path := fmt.Sprintf("%s/%s/servicePrincipals?api-version=1.5", config.GraphURL, creds.TenantID)
	path = path + "&$filter=appId%20eq%20'" + creds.ClientID + "'"
	log.Printf("Get Service Principals request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return "", eh.GenericException(fmt.Sprintf("Error has occurred while parsing params: %v", err))
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", eh.GenericException(fmt.Sprintf("Get Service Principals failed: %s", string(b)))
	}
	var response map[string][]*servicePrincipal

	if err = json.Unmarshal(b, &response); err != nil {
		return "", eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(b)))
	}
	principal := response["value"][0]
	return principal.ObjectID, nil
}

//Assign RBAC role to Application
func assignRoleToApp(c *echo.Context, principalID string, subscription string) error {
	name := uuid.New()
	var properties = map[string]interface{}{
		"properties": map[string]interface{}{
			"roleDefinitionId": fmt.Sprintf("/subscriptions/%s/%s/%s", subscription, authPath, roleContributorID),
			"principalId":      principalID,
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/providers/microsoft.authorization/roleassignments/%s?api-version=%s", config.BaseURL, subscription, name, "2014-10-01-preview")
	log.Printf("Assign RBAC role to Application with params: %s\n", properties)
	log.Printf("Assign RBAC role to Application path: %s\n", path)

	by, err := json.Marshal(properties)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	client, _ := GetAzureClient(c)
	//TODO: handle properly error
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Assign RBAC role to Application failed: %v", err))
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	if response.StatusCode >= 400 {
		return eh.GenericException(fmt.Sprintf("Assign RBAC role to Application failed: %s", string(b)))
	}
	if response.StatusCode != 201 {
		return eh.GenericException(fmt.Sprintf("Assign RBAC role to Application returned status %s with body: %s", response.StatusCode, string(b)))
	}
	return c.NoContent(201)
}
