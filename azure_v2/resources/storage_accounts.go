package resources

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	storageAccountPath = "providers/Microsoft.Storage/storageAccounts"
)

type StorageAccout struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupStorageAccountsRoutes(e *echo.Echo) {
	e.Get("/storage_accounts", listStorageAccounts)
	e.Post("/storage_accounts", createStorageAccount)
	e.Delete("/storage_accounts/:id", deleteStorageAccount)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/storage_accounts")
	// group.Get("", listStorageAccounts)
	// group.Post("", createStorageAccount)
	// group.Delete("/:id", deleteStorageAccount)
}

func listStorageAccounts(c *echo.Context) error {
	return lib.ListResource(c, storageAccountPath, "storage_accounts")
}

func deleteStorageAccount(c *echo.Context) error {
	postParams := c.Request.Form
	group_name := postParams.Get("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, storageAccountPath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

func createStorageAccount(c *echo.Context) error {
	var createParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&createParams)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Group, storageAccountPath, createParams.Name, config.ApiVersion)
	data := StorageAccout{
		Location: createParams.Location,
		Properties: map[string]interface{}{
			"accountType": "Standard_GRS"},
	}

	b, err := lib.CreateResource(c, path, data)
	if err != nil {
		return err
	}
	name := createParams.Name
	// if 'b' is nil we got status 202
	if b != nil {
		var dat *StorageAccout
		if err := json.Unmarshal(b, &dat); err != nil {
			log.Fatal("Unmarshaling failed:", err)
		}
		name = dat.Name
	}

	c.Response.Header().Add("Location", "/storage_accounts/"+name+"?group_name="+createParams.Group)
	return c.NoContent(201)
}
