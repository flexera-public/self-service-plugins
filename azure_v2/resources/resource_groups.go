package resources

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

type ResourceGroup struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty,omitempty"`
	Location   string      `json:"location,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupGroupsRoutes(e *echo.Echo) {
	e.Get("/resource_groups", listResourceGroups)
	e.Post("/resource_groups", createResourceGroup)
	//TODO: add list one action
}

func listResourceGroups(c *echo.Context) error {
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
	groups, err := lib.GetResources(c, path, "/resource_group/%s")
	if err != nil {
		return err
	}
	return c.JSON(200, groups)
}
func createResourceGroup(c *echo.Context) error {
	var createParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
	}

	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&createParams)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	data := ResourceGroup{
		Location: createParams.Location,
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Name, "2015-01-01")
	b, err := lib.CreateResource(c, path, data)
	if err != nil {
		return err
	}
	var dat *ResourceGroup
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	c.Response.Header().Add("Location", "/resource_groups/"+dat.Name)
	return c.NoContent(201)
}
