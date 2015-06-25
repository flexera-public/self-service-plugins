package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

type (
	OperationResponseParams struct {
		OperationId string `json:"operationId"`
		Status      string `json:"status"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime,omitempty"`
		Href        string `json:"href,omitempty"`
	}
	Operation struct {
		Name           string `json:"name,omitempty"`
		Location       string `json:"location,omitempty"`
		ResponseParams OperationResponseParams
	}
)

func SetupOperationRoutes(e *echo.Echo) {
	e.Get("/operations/:id", listOneOperation)
}

func listOneOperation(c *echo.Context) error {
	params := c.Request.Form
	operation := Operation{
		Name:     c.Param("id"),
		Location: params.Get("location"),
	}
	return Get(c, &operation)
}

func (o *Operation) GetResponseParams() interface{} {
	return o.ResponseParams
}

func (o *Operation) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/operations/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, computePath, o.Location, o.Name, "2015-05-01-preview")
}

func (o *Operation) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &o.ResponseParams)
	href := o.GetHref(o.ResponseParams.OperationId, o.Location)
	o.ResponseParams.Href = href
}

func (o *Operation) GetContentType() string {
	return "vnd.rightscale.operation+json"
}

func (o *Operation) GetHref(operationId string, location string) string {
	return fmt.Sprintf("/operations/%s?location=%s", operationId, location)
}

//fake function to support AzureResource by Operation
func (o *Operation) GetCollectionPath(groupName string) string { return "" }

//fake function to support AzureResource by Operation
func (o *Operation) GetRequestParams(c *echo.Context) (interface{}, error) { return "", nil }
