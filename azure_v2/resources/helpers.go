package resources

import (
	"fmt"
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
