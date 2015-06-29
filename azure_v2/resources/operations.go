package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	operationResponseParams struct {
		OperationID string `json:"operationId"`
		Status      string `json:"status"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime,omitempty"`
		Href        string `json:"href,omitempty"`
	}

	// Operation is base struct for Azure Operation resource to store input create params,
	// request create params and response params gotten from cloud.
	Operation struct {
		Name           string `json:"name,omitempty"`
		Location       string `json:"location,omitempty"`
		responseParams operationResponseParams
	}
)

// SetupOperationRoutes declares routes for Operation resource
func SetupOperationRoutes(e *echo.Echo) {
	e.Get("/locations/:location/operations/:id", listOneOperation)
}

func listOneOperation(c *echo.Context) error {
	operation := Operation{
		Name:     c.Param("id"),
		Location: c.Param("location"),
	}
	return Get(c, &operation)
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (o *Operation) GetResponseParams() interface{} {
	return o.responseParams
}

// GetPath returns full path to the sigle operation
func (o *Operation) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/operations/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, computePath, o.Location, o.Name, "2015-05-01-preview")
}

// HandleResponse manage raw cloud response
func (o *Operation) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &o.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	o.responseParams.Href = o.GetHref(o.responseParams.OperationID, o.Location)
	return nil
}

// GetContentType returns content type of operation
func (o *Operation) GetContentType() string {
	return "vnd.rightscale.operation+json"
}

// GetHref returns operation href
func (o *Operation) GetHref(operationID string, location string) string {
	return fmt.Sprintf("/locations/%s/operations/%s", location, operationID)
}

//GetCollectionPath is a fake function to support AzureResource by Operation
func (o *Operation) GetCollectionPath(groupName string) string { return "" }

//GetRequestParams is a fake function to support AzureResource by Operation
func (o *Operation) GetRequestParams(c *echo.Context) (interface{}, error) { return "", nil }
