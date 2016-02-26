package resources

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type operationResponseParams struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
	Href    string `json:"href,omitempty"`
}

// SetupOperationRoutes declares routes for Operation resource
func SetupOperationRoutes(e *echo.Group) {
	e.Get("/locations/:location/services/:service/operations/:id", getOperation)
}

func getOperation(c *echo.Context) error {
	service := c.Param("service")
	var path string
	//Crasy stuff
	if service == "storage" {
		path = fmt.Sprintf("%s/subscriptions/%s/providers/Microsoft.Storage/operations/%s?monitor=true&api-version=%s", config.BaseURL, *config.SubscriptionIDCred, c.Param("id"), "2015-06-15")
	} else if service == "microsoft.compute" {
		path = fmt.Sprintf("%s/subscriptions/%s/providers/Microsoft.Compute/locations/%s/operations/%s?monitor=true&api-version=%s", config.BaseURL, *config.SubscriptionIDCred, c.Param("location"), c.Param("id"), "2015-05-01-preview")
	} else if service == "microsoft.network" {
		path = fmt.Sprintf("%s/subscriptions/%s/providers/Microsoft.Network/locations/%s/operationResults/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, c.Param("location"), c.Param("id"), "2015-06-15")
	} else {
		path = fmt.Sprintf("%s/subscriptions/%s/operationresults/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, c.Param("id"), "2015-11-01")
	}

	client, err := GetAzureClient(c)
	if err != nil {
		return err
	}

	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while requesting resource: %v", err))
	}
	var responseParams operationResponseParams
	responseParams.Href = fmt.Sprintf("locations/%s/operations/%s", c.Param("location"), c.Param("id"))
	if resp.StatusCode == 202 {
		responseParams.Status = "in-progress"
	} else if resp.StatusCode == 200 || resp.StatusCode == 204 {
		responseParams.Status = "succeeded"
	} else {
		details := resp.Header.Get("Location")
		responseParams.Details = fmt.Sprintf("Error has occurred while requesting async operation: %s", details)
		responseParams.Status = "failed"
	}

	return Render(c, 200, responseParams, "vnd.rightscale.operation+json")
}
