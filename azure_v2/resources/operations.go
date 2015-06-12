package resources

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

func SetupOperationRoutes(e *echo.Echo) {
	e.Get("/operations/:id", GetOperation)
}

func GetOperation(c *echo.Context) error {
	params := c.Request.Form
	location := params.Get("location")
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/operations/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, computePath, location, c.Param("id"), "2015-05-01-preview")
	operation, err := lib.GetResource(c, path, "/operations/%s?location="+location)
	if err != nil {
		return err
	}
	return c.JSON(200, operation)
}
