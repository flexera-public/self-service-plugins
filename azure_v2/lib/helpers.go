package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
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

func GetCookie(c *echo.Context, name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", GenericException(fmt.Sprintf("cookie '%s' is missing", cookie))
	}
	return cookie.Value, nil
}

func GetResource(c *echo.Context, path string, href string) (map[string]interface{}, error) {
	client, _ := GetAzureClient(c)
	log.Printf("Get Resource request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return nil, GenericException(fmt.Sprintf("Error has occurred while requesting resource: %v", err))
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 404 {
		return nil, RecordNotFound(c.Param("id"))
	}
	//TODO: add error handling here
	if resp.StatusCode >= 400 {
		return nil, GenericException(fmt.Sprintf("Error has occurred while requesting resource: %s", string(b)))
	}

	var resource map[string]interface{}
	json.Unmarshal(b, &resource)

	resource["href"] = fmt.Sprintf(href, c.Param("id"))

	return resource, nil
}
