package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	RequestParams struct {
		Filter string `json:"filter"`
		Select string `json:"select,omitempty"`
	}
)

// SetupEventsRoutes declares routes for event resource
func SetupEventsRoutes(e *echo.Group) {
	e.Get("/events", listEvents)
}

func listEvents(c *echo.Context) error {
	requestParams := new(RequestParams)
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&requestParams)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	filter := requestParams.Filter
	path := fmt.Sprintf("%s/subscriptions/%s/providers/microsoft.insights/eventtypes/management/values?api-version=%s&$filter=%s", config.BaseURL, *config.SubscriptionIDCred, "2014-04-01", filter)
	if requestParams.Select != "" {
		path = fmt.Sprintf("%s&$select=%s", path, requestParams.Select)
	}

	client, err := GetAzureClient(c)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while getting resource: %v", err))
	}
	req.Header.Add("Content-Type", config.MediaType)
	req.Header.Add("Accept", config.MediaType)
	req.Header.Add("User-Agent", config.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while getting resource: %v", err))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	if resp.StatusCode >= 400 {
		return eh.GenericException(fmt.Sprintf("Error has occurred while getting resource: %s", string(b)))
	}

	var m map[string][]map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		if m["value"] == nil {
			return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(b)))
		}
	}

	//TODO: add hrefs or use AzureResource interface
	return c.JSON(200, m["value"])
}
