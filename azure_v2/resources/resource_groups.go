package resources

import (
	"fmt"

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
}

func listResourceGroups(c *echo.Context) error {
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
	groups, err := lib.GetResources(c, path, "/azure_plugin/resource_group/%s")
	if err != nil {
		return err
	}
	return c.JSON(200, groups)
}
