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
	e.Delete("/networks/:id", deleteNetwork)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/networks")
	// group.Get("", listNetworks)
	// group.Post("", createNetwork)
	// group.Delete("/:id", deleteNetwork)
}

func listNetworks(c *echo.Context) error {
	return lib.ListResource(c, networkPath, "networks")
}

func deleteNetwork(c *echo.Context) error {
	postParams := c.Request.Form
	group_name := postParams.Get("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, networkPath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

func createNetwork(c *echo.Context) error {
	var createParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&createParams)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Group, networkPath, createParams.Name, config.ApiVersion)
	var subnets []map[string]interface{}
	data := Network{
		Name:     createParams.Name,
		Location: createParams.Location,
		Properties: map[string]interface{}{
			"addressSpace": map[string]interface{}{
				"addressPrefixes": []string{"10.0.0.0/16"},
			},
			"subnets": append(subnets, map[string]interface{}{
				"name": createParams.Name,
				"properties": map[string]interface{}{
					"addressPrefix": "10.0.0.0/16",
				},
			}),
		},
	}

	b, err := lib.CreateResource(c, path, data)
	if err != nil {
		return err
	}
	var dat *Network
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}
	c.Response.Header().Add("Location", "/networks/"+dat.Name+"?group_name="+createParams.Group)
	return c.NoContent(201)
}
