package lib

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

// Retrieve client initialized by middleware, send error response if not found
// This function should be used by controller actions that need to use the client
func GetAzureClient(c *echo.Context) (*http.Client, error) {
	client, _ := c.Get("azure").(*http.Client)
	if client == nil {
		return nil, GenericException(fmt.Sprintf("failed to retrieve Azure client, check middleware"))
	}
	return client, nil
}

func GetCookie(c *echo.Context, name string) (*http.Cookie, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return nil, GenericException(fmt.Sprintf("cookie '%s' is missing", cookie))
	}
	return cookie, nil
}

func ListNestedResources(c *echo.Context, parentPath string, relativePath string) ([]interface{}, error) {
	re := regexp.MustCompile("\\?(.+)") //remove tail with uri params
	parentResources, err := GetResources(c, parentPath)
	if err != nil {
		return nil, err
	}

	var resources []interface{}
	for _, parent := range parentResources {
		resource := parent.(map[string]interface{})
		resourcePath := re.ReplaceAllLiteralString(parentPath, "/") + resource["name"].(string) + relativePath
		resp, err := GetResources(c, resourcePath)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resp...)
	}
	return resources, nil
}

func ListResource(c *echo.Context, resourcePath string) error {
	group_name := c.Param("group_name")
	if group_name != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, resourcePath, config.ApiVersion)
		resources, err := GetResources(c, path)
		if err != nil {
			return err
		}
		return c.JSON(200, resources)
	} else {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
		relativePath := fmt.Sprintf("/%s?api-version=%s", resourcePath, config.ApiVersion)
		resources, err := ListNestedResources(c, path, relativePath)
		if err != nil {
			return err
		}
		return c.JSON(200, resources)
	}

}

func GetResources(c *echo.Context, path string) ([]interface{}, error) {
	client, _ := GetAzureClient(c)
	log.Printf("Get Resources request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return nil, GenericException(fmt.Sprintf("Error has occurred while requesting resources: %v", err))
	}
	b, _ := ioutil.ReadAll(resp.Body)
	//TODO: add error handling here
	if resp.StatusCode >= 400 {
		return nil, GenericException(fmt.Sprintf("Error has occurred while requesting resources: %s", string(b)))
	}

	var m map[string][]interface{}
	json.Unmarshal(b, &m)

	return m["value"], nil
}

func DeleteResource(c *echo.Context, path string) error {
	client, _ := GetAzureClient(c)
	log.Printf("Delete request: %s\n", path)

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return GenericException(fmt.Sprintf("Error has occurred while deleting resource: %v", err))
	}

	resp, err := client.Do(req)

	if err != nil {
		return GenericException(fmt.Sprintf("Error has occurred while deleting resource: %v", err))
	}

	return c.JSON(resp.StatusCode, "")
}
