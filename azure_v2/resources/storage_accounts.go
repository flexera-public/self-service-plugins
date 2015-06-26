package resources

import (
	"encoding/json"
	"fmt"

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
func SetupStorageAccountsRoutes(e *echo.Echo) {
	e.Get("/storage_accounts", listStorageAccounts)
	e.Get("/storage_accounts/:id", listOneStorageAccount)
	e.Post("/storage_accounts", createStorageAccount)
	e.Delete("/storage_accounts/:id", deleteStorageAccount)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/storage_accounts")
	// group.Get("", listStorageAccounts)
	// group.Post("", createStorageAccount)
	// group.Delete("/:id", deleteStorageAccount)
}

func listStorageAccounts(c *echo.Context) error {
	return List(c, new(StorageAccount))
}

func listOneStorageAccount(c *echo.Context) error {
	params := c.Request.Form
	storageAccount := StorageAccount{
		createParams: storageAccountCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Get(c, &storageAccount)
}

func createStorageAccount(c *echo.Context) error {
	storageAccount := new(StorageAccount)
	return Create(c, storageAccount)
}

func deleteStorageAccount(c *echo.Context) error {
	params := c.Request.Form
	storageAccount := StorageAccount{
		createParams: storageAccountCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
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

	s.requestParams.Location = s.createParams.Location
	s.requestParams.Properties = map[string]interface{}{"accountType": "Standard_GRS"}

	return s.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (s *StorageAccount) GetResponseParams() interface{} {
	return s.responseParams
}

// GetPath returns full path to the sigle storage account
func (s *StorageAccount) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, s.createParams.Group, storageAccountPath, s.createParams.Name, config.APIVersion)
}

// GetCollectionPath returns full path to the collection of storage accounts
func (s *StorageAccount) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, storageAccountPath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (s *StorageAccount) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &s.responseParams)
	href := s.GetHref(s.createParams.Group, s.responseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		s.responseParams.Href = href
	}
}

// GetContentType returns storage account content type
func (s *StorageAccount) GetContentType() string {
	return "vnd.rightscale.storage_account+json"
}

// GetHref returns storage account href
func (s *StorageAccount) GetHref(groupName string, networkName string) string {
	return fmt.Sprintf("/storage_accounts/%s?group_name=%s", networkName, groupName)
}
