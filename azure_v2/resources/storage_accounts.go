package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	storageAccountPath = "providers/Microsoft.Storage/storageAccounts"
)

type (
	storageAccountResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	storageAccountRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	storageAccountCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	// StorageAccount is base struct for Azure Storage Account resource to store input create params,
	// request create params and response params gotten from cloud.
	StorageAccount struct {
		createParams   storageAccountCreateParams
		requestParams  storageAccountRequestParams
		responseParams storageAccountResponseParams
	}
)

// SetupStorageAccountsRoutes declares routes for Storage account resource
func SetupStorageAccountsRoutes(e *echo.Group) {
	e.Get("/storage_accounts", listStorageAccounts)
	// e.Get("/storage_accounts/:id", listOneStorageAccount)
	// e.Post("/storage_accounts", createStorageAccount)
	// e.Delete("/storage_accounts/:id", deleteStorageAccount)

	//nested routes
	group := e.Group("/resource_groups/:group_name/storage_accounts")
	group.Get("", listStorageAccounts)
	group.Get("/:name", listOneStorageAccount)
	group.Get("/:name/check_name", checkNameAvailability)
	group.Post("", createStorageAccount)
	group.Delete("/:name", deleteStorageAccount)
	//group.Delete("/:id/keys", getStorageAccountKeys)
}

func listStorageAccounts(c *echo.Context) error {
	return List(c, new(StorageAccount))
}

func listOneStorageAccount(c *echo.Context) error {
	storageAccount := StorageAccount{
		createParams: storageAccountCreateParams{
			Name:  c.Param("name"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &storageAccount)
}

func createStorageAccount(c *echo.Context) error {
	storageAccount := new(StorageAccount)
	return Create(c, storageAccount)
}

func deleteStorageAccount(c *echo.Context) error {
	storageAccount := StorageAccount{
		createParams: storageAccountCreateParams{
			Name:  c.Param("name"),
			Group: c.Param("group_name"),
		},
	}
	return Delete(c, &storageAccount)
}

// GetRequestParams prepares parameters for create storage account request to the cloud
func (s *StorageAccount) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&s.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	s.createParams.Group = c.Param("group_name")

	s.requestParams.Location = s.createParams.Location
	s.requestParams.Properties = map[string]interface{}{"accountType": "Standard-RAGRS"}

	return s.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (s *StorageAccount) GetResponseParams() interface{} {
	return s.responseParams
}

// GetPath returns full path to the sigle storage account
func (s *StorageAccount) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, s.createParams.Group, storageAccountPath, s.createParams.Name, "2015-06-15")
}

// GetCollectionPath returns full path to the collection of storage accounts
func (s *StorageAccount) GetCollectionPath(groupName string) string {
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, storageAccountPath, "2015-06-15")
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, storageAccountPath, "2015-06-15")
}

// HandleResponse manage raw cloud response
func (s *StorageAccount) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &s.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	if actionName == "create" {
		c.Response().Header().Add("Location", fmt.Sprintf("resource_groups/%s/storage_accounts/%s", s.createParams.Group, s.createParams.Name))
	} else if actionName == "get" {
		s.responseParams.Href = s.GetHref(s.responseParams.ID)
	}
	return nil
}

// GetContentType returns storage account content type
func (s *StorageAccount) GetContentType() string {
	return "vnd.rightscale.storage_account+json"
}

// GetHref returns storage account href
func (s *StorageAccount) GetHref(accountID string) string {
	array := strings.Split(accountID, "/")
	return fmt.Sprintf("resource_groups/%s/storage_accounts/%s", array[len(array)-5], array[len(array)-1])
}

func checkNameAvailability(c *echo.Context) error {
	client, err := GetAzureClient(c)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/subscriptions/%s/providers/Microsoft.Storage/checkNameAvailability?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, "2015-06-15")
	by, err := json.Marshal(map[string]interface{}{
		"name": c.Param("name"),
		"type": "Microsoft.Storage/storageAccounts",
	})
	if err != nil {
		eh.GenericException(fmt.Sprintf("Error has occurred while marshaling data: %v", err))
	}
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	req, err := http.NewRequest("POST", path, reader)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while registering provider: %v", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	return Render(c, 200, response, "application/json")

}
