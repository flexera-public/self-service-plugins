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
	e.Post("/application/register", assignRoleToApp)
	e.Delete("/application/unregister", unassignRoleFromApp)
}

//Assign RBAC role to Application
func assignRoleToApp(c *echo.Context) error {
	principalID, subscription, err := prepareParams(c)
	if err != nil {
		return err
	}
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
	if err != nil {
		eh.GenericException(fmt.Sprintf("Error has occurred while marshaling data: %v", err))
	}
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	client, err := GetAzureClient(c)
	if err != nil {
		return err
	}
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

// Delete Role assignment in order to un-register application
func unassignRoleFromApp(c *echo.Context) error {
	principalID, subscription, err := prepareParams(c)
	if err != nil {
		return err
	}
	name, err := findRoleAssignmentID(c, principalID, subscription)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/subscriptions/%s/providers/microsoft.authorization/roleassignments/%s?api-version=%s", config.BaseURL, subscription, name, "2014-10-01-preview")
	log.Printf("Unassign RBAC role from Application path: %s\n", path)

	req, err := http.NewRequest("DELETE", path, nil)
	req.Header.Add("User-Agent", config.UserAgent)
	client, err := GetAzureClient(c)
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Unassignment RBAC role from Application failed: %v", err))
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	if response.StatusCode >= 400 {
		return eh.GenericException(fmt.Sprintf("Unassignment RBAC role from Application failed: %s", string(b)))
	}
	return c.NoContent(204)
}

func findRoleAssignmentID(c *echo.Context, principalID string, subscription string) (string, error) {
	roleDefinitionID := fmt.Sprintf("/subscriptions/%s/%s/%s", subscription, authPath, roleContributorID)
	path := fmt.Sprintf("%s/subscriptions/%s/providers/microsoft.authorization/roleassignments?api-version=%s", config.BaseURL, subscription, "2014-10-01-preview")
	roleAssignments, err := GetResources(c, path)
	if err != nil {
		return "", err
	}

	for _, ra := range roleAssignments {
		// find assignment for gotten roleDefinitionID and principalID
		if ra["properties"].(map[string]interface{})["roleDefinitionId"].(string) == roleDefinitionID && ra["properties"].(map[string]interface{})["principalId"].(string) == principalID {
			return ra["name"].(string), nil
		}
	}
	return "", eh.GenericException(fmt.Sprintf("Role assignment is not found for role definition '%s' and principal ID '%s'.", roleDefinitionID, principalID))
}

// Get params required for app (un)registration
func prepareParams(c *echo.Context) (string, string, error) {
	creds := am.Credentials{
		GrantType:    "client_credentials",
		Resource:     "https://graph.windows.net/",
		TenantID:     *config.TenantIDCred,
		ClientID:     *config.ClientIDCred,
		ClientSecret: *config.ClientSecretCred,
		RefreshToken: *config.RefreshTokenCred,
		Subscription: *config.SubscriptionIDCred,
	}
	authResponse, err := creds.RequestToken()
	if err != nil {
		return "", "", err
	}
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: authResponse.AccessToken}}
	graphClient := t.Client()
	principalID, err := getServicePrincipal(graphClient, &creds)
	if err != nil {
		return "", "", err
	}
	return principalID, creds.Subscription, nil
}

func getServicePrincipal(client *http.Client, creds *am.Credentials) (string, error) {
	path := fmt.Sprintf("%s/%s/servicePrincipals?api-version=1.5", config.GraphURL, creds.TenantID)
	path = path + "&$filter=appId%20eq%20'" + creds.ClientID + "'"
	log.Printf("Get Service Principals request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return "", eh.GenericException(fmt.Sprintf("Error has occurred while sending request: %v", err))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	if resp.StatusCode >= 400 {
		return "", eh.GenericException(fmt.Sprintf("Get Service Principals failed: %s", string(b)))
	}
	var response map[string][]*servicePrincipal

	if err = json.Unmarshal(b, &response); err != nil {
		if response["value"] == nil {
			// return resources if unmarshaling is success for 'value' key
			// error occurs if value of the hash is not a []map[string]interface{}
			return "", eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(b)))
		}
	}

	principal := response["value"][0]
	return principal.ObjectID, nil
}
