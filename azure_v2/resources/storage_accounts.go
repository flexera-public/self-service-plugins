package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

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

	//nested routes
	group := e.Group("/resource_groups/:group_name/storage_accounts")
	group.Get("", listStorageAccounts)
	// group.Post("", createInstance)
	// group.Delete("/:id", deleteInstance)
}

func listStorageAccounts(c *echo.Context) error {
	return lib.ListResource(c, storageAccountPath)
}

func createStorageAccount(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), storageAccountPath, postParams.Get("name"), config.ApiVersion)
	log.Printf("Create Storage Account request with params: %s\n", postParams)
	log.Printf("Create Storage Account path: %s\n", path)
	data := StorageAccout{
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"accountType": "Standard_GRS"},
	}

	by, err := json.Marshal(data)
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
	b, _ := ioutil.ReadAll(response.Body)
	//TODO: handle 202 state
	return c.JSON(response.StatusCode, string(b))
}
