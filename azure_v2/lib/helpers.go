package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

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

	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return GenericException(fmt.Sprintf("Error has occurred while deleting resource: %v", string(b)))
	}

	return c.JSON(204, "")
}

func CreateResource(c *echo.Context, path string, createParams interface{}) ([]byte, error) {
	client, _ := GetAzureClient(c)

	log.Printf("Create Resource request with params: %#v\n", createParams)
	log.Printf("Create Resource path: %s\n", path)

	by, err := json.Marshal(createParams)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, err := http.NewRequest("PUT", path, reader)
	if err != nil {
		return nil, GenericException(fmt.Sprintf("Error has occurred while creating resource: %v", err))
	}
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		return nil, GenericException(fmt.Sprintf("Error has occurred while creating resource: %v", err))
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return nil, GenericException(fmt.Sprintf("Error has occurred while creating resource: %s", string(b)))
	}
	if response.Header.Get("azure-asyncoperation") != "" {
		c.Response.Header().Add("azure-asyncoperation", response.Header.Get("azure-asyncoperation"))
	}
	if response.StatusCode == 202 {
		return nil, nil
	}
	return b, nil
}
