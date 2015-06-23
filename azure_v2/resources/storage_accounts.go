package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	storageAccountPath = "providers/Microsoft.Storage/storageAccounts"
)

type (
	StorageAccountResponseParams struct {
		Id         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	StorageAccountRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	StorageAccountCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	StorageAccount struct {
		CreateParams   StorageAccountCreateParams
		RequestParams  StorageAccountRequestParams
		ResponseParams StorageAccountResponseParams
	}
)

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
	return lib.List(c, new(StorageAccount))
}

func listOneStorageAccount(c *echo.Context) error {
	params := c.Request.Form
	storage_account := StorageAccount{
		CreateParams: StorageAccountCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return lib.Get(c, &storage_account)
}

func createStorageAccount(c *echo.Context) error {
	storage_account := new(StorageAccount)
	return lib.Create(c, storage_account)
}

func deleteStorageAccount(c *echo.Context) error {
	params := c.Request.Form
	storage_account := StorageAccount{
		CreateParams: StorageAccountCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return lib.Delete(c, &storage_account)
}

func (s *StorageAccount) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&s.CreateParams)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	s.RequestParams.Location = s.CreateParams.Location
	s.RequestParams.Properties = map[string]interface{}{"accountType": "Standard_GRS"}

	return s.RequestParams, nil
}

func (s *StorageAccount) GetResponseParams() interface{} {
	return s.ResponseParams
}

func (s *StorageAccount) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, s.CreateParams.Group, storageAccountPath, s.CreateParams.Name, config.ApiVersion)
}

func (s *StorageAccount) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, storageAccountPath, config.ApiVersion)
}

func (s *StorageAccount) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &s.ResponseParams)
	href := s.GetHref(s.CreateParams.Group, s.ResponseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		s.ResponseParams.Href = href
	}
}

func (s *StorageAccount) GetContentType() string {
	return "vnd.rightscale.storage_account+json"
}

func (s *StorageAccount) GetHref(groupName string, networkName string) string {
	return fmt.Sprintf("/storage_accounts/%s?group_name=%s", groupName, networkName)
}
