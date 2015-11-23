package resources

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

// SetupGroupsRoutes declares routes for resource group resource
func SetupInstanceTypesRoutes(e *echo.Group) {
	e.Get("/locations/:location/instance_types", listInstanceTypes)
}

//This API lists all available virtual machine sizes for a subscription in a given region.
func listInstanceTypes(c *echo.Context) error {
	location := c.Param("location")
	path := fmt.Sprintf("%s/subscriptions/%s/providers/Microsoft.Compute/locations/%s/vmSizes?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, location, "2015-06-15")
	its, err := GetResources(c, path)
	if err != nil {
		return err
	}

	//TODO: add hrefs or use AzureResource interface
	return c.JSON(200, its)
}
