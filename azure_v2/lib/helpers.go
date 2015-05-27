package lib

import (
	// "bytes"
	"encoding/json"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	//"github.com/rightscale/self-service-plugins/azure_v2/config"
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

// func ListNestedResources(c *echo.Context, relativePath string) error {
// 	code, resp := GetResources(c, "")

// 	var resources []interface{}
// 	for _, resource_group := range resp {
// 		name := resource_group.(*ResourceGroup).Name
// 		_, resp := GetResources(c, "/"+name+"/"+relativePath)
// 		resources = append(resources, resp...)
// 	}
// 	// [].to_json => null ... why?
// 	return c.JSON(code, resources)
// }

func GetResources(c *echo.Context, path string) error {
	client, _ := GetAzureClient(c)
	log.Printf("Get Resources request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return GenericException(fmt.Sprintf("Error has occurred while requesting resources: %v", err))
	}
	b, _ := ioutil.ReadAll(resp.Body)
	//TODO: add error handling here
	if resp.StatusCode >= 400 {
		return GenericException(fmt.Sprintf("Error has occurred while requesting resources: %s", string(b)))
	}

	var m map[string][]interface{}
	json.Unmarshal(b, &m)

	return c.JSON(resp.StatusCode, m["value"])
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
