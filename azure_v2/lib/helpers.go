package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func GetCookie(c *echo.Context, name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", GenericException(fmt.Sprintf("cookie '%s' is missing", cookie))
	}
	return cookie.Value, nil
}

// JSON sends an resource specific content type response with status code.
func Render(c *echo.Context, code int, i interface{}, contentType string) error {
	c.Response.Header().Set(echo.ContentType, contentType) //ApplicationJSON
	c.Response.WriteHeader(code)
	return json.NewEncoder(c.Response).Encode(i)
}

func ListNestedResources(c *echo.Context, parentPath string, relativePath string, resourcesName string) ([]map[string]interface{}, error) {
	re := regexp.MustCompile("\\?(.+)") //remove tail with uri params
	parentResources, err := GetResources(c, parentPath, "/resource_group/%s")
	if err != nil {
		return nil, err
	}

	var resources []map[string]interface{}
	for _, parent := range parentResources {
		resourcePath := re.ReplaceAllLiteralString(parentPath, "/") + parent["name"].(string) + relativePath
		resp, err := GetResources(c, resourcePath, "/"+resourcesName+"/%s?group_name="+parent["name"].(string))
		if err != nil {
			return nil, err
		}
		resources = append(resources, resp...)
	}
	return resources, nil
}

func ListResource(c *echo.Context, resourcePath string, resourcesName string) error {
	group_name := c.Param("group_name")
	var resources []map[string]interface{}
	var err error
	if group_name != "" {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, resourcePath, config.ApiVersion)
		resources, err = GetResources(c, path, "/"+resourcesName+"/%s?group_name="+group_name)
		if err != nil {
			return err
		}
	} else {
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
		relativePath := fmt.Sprintf("/%s?api-version=%s", resourcePath, config.ApiVersion)
		resources, err = ListNestedResources(c, path, relativePath, resourcesName)
		if err != nil {
			return err
		}
	}
	// init empty array
	if len(resources) == 0 {
		resources = make([]map[string]interface{}, 0)
	}
	return Render(c, 200, resources, "vnd.rightscale."+resourcesName+"+json")

}

func GetResources(c *echo.Context, path string, href string) ([]map[string]interface{}, error) {
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

	var m map[string][]map[string]interface{}
	var resources []map[string]interface{}
	err = json.Unmarshal(b, &m)
	resources = m["value"]
	if err != nil {
		//try to unmarshal with different interface
		json.Unmarshal(b, &resources)
	}

	//add href for each resource
	for _, resource := range resources {
		resource["href"] = fmt.Sprintf(href, resource["name"])
	}

	return resources, nil
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

	log.Printf("Create Instances request with params: %#v\n", createParams)
	log.Printf("Create Instances path: %s\n", path)

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
