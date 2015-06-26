package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

// GetAzureClient retrieves client initialized by middleware, send error response if not found
// This function should be used by controller actions that need to use the client
func GetAzureClient(c *echo.Context) (*http.Client, error) {
	client, _ := c.Get("azure").(*http.Client)
	if client == nil {
		return nil, eh.GenericException(fmt.Sprintf("failed to retrieve Azure client, check middleware"))
	}
	return client, nil
}

// GetResource sends requests to the clouds to get resource
func GetResource(c *echo.Context, path string, href string) (map[string]interface{}, error) {
	client, err := GetAzureClient(c)
	if err != nil {
		return nil, err
	}
	log.Printf("Get Resource request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while requesting resource: %v", err))
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	if resp.StatusCode == 404 {
		return nil, eh.RecordNotFound(c.Param("id"))
	}
	if resp.StatusCode >= 400 {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while requesting resource: %s", string(b)))
	}

	var resource map[string]interface{}
	if err = json.Unmarshal(b, &resource); err != nil {
		return nil, eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(b)))
	}

	resource["href"] = fmt.Sprintf(href, c.Param("id"))

	return resource, nil
}
