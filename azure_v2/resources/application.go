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

func SetupAuthRoutes(e *echo.Echo) {
	e.Post("/application/register", registerApp)
}

// Get App-Only Access Token for Azure AD Graph API
func registerApp(c *echo.Context) error {
	var registrationCreds struct {
		TenantId     string `json:"tenant"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Subscription string `json:"subscription"`
	}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&registrationCreds)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	authResponse, err := lib.RequestToken(registrationCreds.TenantId, "client_credentials", "https://graph.windows.net/", registrationCreds.ClientId, registrationCreds.ClientSecret, "")
	if err != nil {
		return err
	}
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: authResponse.AccessToken}}
	graphClient := t.Client()
	principalId, err1 := getServicePrincipal(graphClient, registrationCreds.TenantId, registrationCreds.ClientId)
	if err1 != nil {
		return err1
	}
	return assignRoleToApp(c, principalId, registrationCreds.Subscription)
}

// {"odata.metadata":"https://graph.windows.net/09b8fec1-4b8d-48dd-8afa-5c1a775ea0f2/$metadata#directoryObjects/Microsoft.DirectoryServices.ServicePrincipal",
//  "value":[{
// 		"odata.type":"Microsoft.DirectoryServices.ServicePrincipal",
// 		"objectType":"ServicePrincipal",
// 		"objectId":"7f06b355-4136-4d8b-a3c8-7028f59869ae",
// 		"deletionTimestamp":null,
// 		"accountEnabled":true,
// 		"appDisplayName":"RightScale",
// 		"appId":"57e3974e-7bb8-47a3-8b28-1f4026b6ac65",
// 		"appOwnerTenantId":"b6bf44a6-1833-4f1e-aff2-52c7a6834b63",
// 		"appRoleAssignmentRequired":false,
// 		"appRoles":[],
// 		"displayName":"RightScale",
// 		"errorUrl":null,
//	 	"homepage":"https://my.rightscale.com",
// 		"keyCredentials":[],
// 		"logoutUrl":null,
// 		"oauth2Permissions":[{"adminConsentDescription":"Allow the application to access RightScale.com on behalf of the signed-in user.","adminConsentDisplayName":"Access RightScale.com","id":"677c5e7a-82f1-4889-8ff8-21b61569c8c0","isEnabled":true,"type":"User","userConsentDescription":"Allow the application to access RightScale.com on your behalf.","userConsentDisplayName":"Access RightScale.com","value":"user_impersonation"}],
// 		"passwordCredentials":[],
// 		"preferredTokenSigningKeyThumbprint":null,
// 		"publisherName":"RightScale Test",
// 		"replyUrls":["https://ad.test.rightscale.com"],
// 		"samlMetadataUrl":null,
// 		"servicePrincipalNames":["https://test.rightscale.com","57e3974e-7bb8-47a3-8b28-1f4026b6ac65"],
//  "tags":["WindowsAzureActiveDirectoryIntegratedApp"]}]}
func getServicePrincipal(client *http.Client, tenantId string, cliectId string) (string, error) {
	path := fmt.Sprintf("%s/%s/servicePrincipals?api-version=1.5", config.GraphUrl, tenantId)
	path = path + "&$filter=appId%20eq%20'" + cliectId + "'"
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
	var response map[string]interface{}

	if err = json.Unmarshal(b, &response); err != nil {
		return "", lib.GenericException(fmt.Sprintf("got bad response from server: %s", string(b)))
	}
	principal := response["value"].([]interface{})[0].(map[string]interface{})
	return principal["objectId"].(string), nil
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
