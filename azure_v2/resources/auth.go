package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/goauth2/oauth"
	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	authPath    = "providers/Microsoft.Authorization/roleDefinitions"
	roleOwnerId = "8e3af657-a8ff-443c-a75c-2fe8c4bcb635"
)

//subscriptions/09cbd307-aa71-4aca-b346-5f253e6e3ebb/providers/Microsoft.Authorization/roleDefinitions?api-version=2014-07-01-preview
func SetupAuthRoutes(e *echo.Echo) {
	e.Get("/auth_app", getAppToken)
}

//https://graph.windows.net/62e173e9-301e-423e-bcd4-29121ec1aa24/servicePrincipals?api-version=1.5&$filter=appId%20eq%20'a0448380-c346-4f9f-b897-c18733de9394'
func getServicePrincipals(client *http.Client) (*http.Response, error) {
	path := fmt.Sprintf("https://graph.windows.net/%s/servicePrincipals?api-version=1.5", *config.TenantIdCred)
	path = path + "&$filter=appId%20eq%20'" + *config.ClientIdCred + "'"
	log.Printf("Get Service Principals request: %s\n", path)
	resp, err := client.Get(path)

	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get App-Only Access Token for Azure AD Graph API
func getAppToken(c *echo.Context) error {
	resp, err := lib.RequestToken("client_credentials", "https://graph.windows.net/")
	if err != nil {
		log.Fatal("POST:", err)
	}

	t := &oauth.Transport{Token: &oauth.Token{AccessToken: resp.AccessToken}}
	client := t.Client()
	principalsResponse, err1 := getServicePrincipals(client)
	if err1 != nil {
		log.Fatal("GET:", err)
	}
	b, _ := ioutil.ReadAll(principalsResponse.Body)

	return c.JSON(principalsResponse.StatusCode, string(b))
}

//Assign RBAC role to Application
func assignRoleToApp(client *http.Client, principalId string) (*http.Response, error) {
	// PUT https://management.azure.com/subscriptions/09cbd307-aa71-4aca-b346-5f253e6e3ebb/providers/microsoft.authorization/roleassignments/4f87261d-2816-465d-8311-70a27558df4c?api-version=2014-10-01-preview HTTP/1.1
	// Authorization: Bearer eyJ0eXAiOiJKV1QiL*****FlwO1mM7Cw6JWtfY2lGc5A
	// Content-Type: application/json
	// Content-Length: 230

	// {"properties": {"roleDefinitionId":"/subscriptions/09cbd307-aa71-4aca-b346-5f253e6e3ebb/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7","principalId":"c3097b31-7309-4c59-b4e3-770f8406bad2"}}

	name := config.TenantIdCred //is a new guid created for the new role assignment...let's use tenant id for now
	var properties = map[string]interface{}{
		"properties": map[string]interface{}{
			"roleDefinitionId": fmt.Sprintf("/subscriptions/%s/%s/%s", *config.SubscriptionIdCred, authPath, name),
			"principalId":      principalId,
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/providers/microsoft.authorization/roleassignments/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, name, config.ApiVersion)
	log.Printf("Assign RBAC role to Application with params: %s\n", properties)
	log.Printf("Assign RBAC role to Application path: %s\n", path)

	by, err := json.Marshal(properties)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	log.Printf("READER: %s", reader)
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("PUT:", err)
	}
	defer response.Body.Close()
	return response, nil
}
