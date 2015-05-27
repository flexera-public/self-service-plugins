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
	networkPath = "providers/Microsoft.Network/virtualNetworks"
)

type Network struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Tags       interface{} `json:"tags,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupNetworkRoutes(e *echo.Echo) {
	e.Get("/networks", listNetworks)
	e.Post("/networks", createNetwork)

	//nested routes
	group := e.Group("/resource_groups/:group_name/networks")
	group.Get("", listNetworks)
	// group.Post("", createInstance)
	// group.Delete("/:id", deleteInstance)
}

func listNetworks(c *echo.Context) error {
	return lib.ListResource(c, networkPath)
}

func createNetwork(c *echo.Context) error {
	postParams := c.Request.Form
	client, _ := lib.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), networkPath, postParams.Get("name"), config.ApiVersion)
	log.Printf("Create Network request with params: %s\n", postParams)
	log.Printf("Create Network path: %s\n", path)
	var subnets []map[string]interface{}
	data := Network{
		Name:     postParams.Get("name"),
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"addressSpace": map[string]interface{}{
				"addressPrefixes": []string{"10.0.0.0/16"},
			},
			"subnets": append(subnets, map[string]interface{}{
				"name": postParams.Get("name"),
				"properties": map[string]interface{}{
					"addressPrefix": "10.0.0.0/16",
				},
			}),
		},
	}

	by, err := json.Marshal(data)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
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
	var dat *Network
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	return c.JSON(response.StatusCode, dat)
}
