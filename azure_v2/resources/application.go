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
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	authPath          = "providers/Microsoft.Authorization/roleDefinitions"
	roleContributorId = "b24988ac-6180-42a0-ab88-20f7382dd24c"
)

type ServicePrincipal struct {
	// Add more fields if needed
	ObjectId string `json:"objectId"`
}

func SetupAuthRoutes(e *echo.Echo) {
	e.Post("/application/register", registerApp)
}

// Get App-Only Access Token for Azure AD Graph API
func registerApp(c *echo.Context) error {
	creds := lib.Credentials{GrantType: "client_credentials", Resource: "https://graph.windows.net/"}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&creds)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	authResponse, err := creds.RequestToken()
	if err != nil {
		return err
	}
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: authResponse.AccessToken}}
	graphClient := t.Client()
	principalId, err := getServicePrincipal(graphClient, &creds)
	if err != nil {
		return err
	}
	return assignRoleToApp(c, principalId, creds.Subscription)
}

func getServicePrincipal(client *http.Client, creds *lib.Credentials) (string, error) {
	path := fmt.Sprintf("%s/%s/servicePrincipals?api-version=1.5", config.GraphUrl, creds.TenantId)
	path = path + "&$filter=appId%20eq%20'" + creds.ClientId + "'"
	log.Printf("Get Service Principals request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return "", lib.GenericException(fmt.Sprintf("Error has occurred while parsing params: %v", err))
	}

	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", lib.GenericException(fmt.Sprintf("Get Service Principals failed: %s", string(b)))
	}
	var response map[string][]*ServicePrincipal

	if err = json.Unmarshal(b, &response); err != nil {
		return "", lib.GenericException(fmt.Sprintf("got bad response from server: %s", string(b)))
	}
	principal := response["value"][0]
	return principal.ObjectId, nil
}

//Assign RBAC role to Application
func assignRoleToApp(c *echo.Context, principalId string, subscription string) error {
	name := uuid.New()
	var properties = map[string]interface{}{
		"properties": map[string]interface{}{
			"roleDefinitionId": fmt.Sprintf("/subscriptions/%s/%s/%s", subscription, authPath, roleContributorId),
			"principalId":      principalId,
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/providers/microsoft.authorization/roleassignments/%s?api-version=%s", config.BaseUrl, subscription, name, "2014-10-01-preview")
	log.Printf("Assign RBAC role to Application with params: %s\n", properties)
	log.Printf("Assign RBAC role to Application path: %s\n", path)

	by, err := json.Marshal(properties)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	client, _ := lib.GetAzureClient(c)
	//TODO: handle properly error
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Assign RBAC role to Application failed: %v", err))
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	if response.StatusCode >= 400 {
		return lib.GenericException(fmt.Sprintf("Assign RBAC role to Application failed: %s", string(b)))
	}
	if response.StatusCode != 201 {
		return lib.GenericException(fmt.Sprintf("Assign RBAC role to Application returned status %s with body: %s", response.StatusCode, string(b)))
	}
	return c.NoContent(201)
}
