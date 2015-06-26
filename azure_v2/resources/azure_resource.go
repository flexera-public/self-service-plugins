package resources

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
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

// AzureResource is interface which should support every resource in order to use generic functions List/Get/Create/Delete
type AzureResource interface {
	GetRequestParams(*echo.Context) (interface{}, error)
	GetResponseParams() interface{}
	GetPath() string
	GetCollectionPath(string) string
	HandleResponse(*echo.Context, []byte, string)
	GetContentType() string
	GetHref(string, string) string
}

// Create new resource
func Create(c *echo.Context, r AzureResource) error {
	client, _ := GetAzureClient(c)
	requestParams, _ := r.GetRequestParams(c)

	by, err := json.Marshal(requestParams)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	path := r.GetPath()
	request, err := http.NewRequest("PUT", path, reader)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while creating resource: %v", err))
	}
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while creating resource: %v", err))
	}
	//r.HandleResponse(response)
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return eh.GenericException(fmt.Sprintf("Error has occurred while creating resource: %s", string(b)))
	}
	if response.Header.Get("azure-asyncoperation") != "" {
		c.Response.Header().Add("azure-asyncoperation", response.Header.Get("azure-asyncoperation"))
	}

	r.HandleResponse(c, b, "create")
	return c.NoContent(201)
}

// Delete resource
func Delete(c *echo.Context, r AzureResource) error {
	client, _ := GetAzureClient(c)
	path := r.GetPath()
	log.Printf("Delete request: %s\n", path)

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while deleting resource: %v", err))
	}

	resp, err := client.Do(req)

	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while deleting resource: %v", err))
	}

	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return eh.GenericException(fmt.Sprintf("Error has occurred while deleting resource: %v", string(b)))
	}

	return c.JSON(204, "")
}

// Get resource
func Get(c *echo.Context, r AzureResource) error {
	client, _ := GetAzureClient(c)
	path := r.GetPath()
	log.Printf("Get Resource request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while requesting resource: %v", err))
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 404 {
		return eh.RecordNotFound(c.Param("id"))
	}
	if resp.StatusCode >= 400 {
		return eh.GenericException(fmt.Sprintf("Error has occurred while requesting resource: %s", string(b)))
	}

	r.HandleResponse(c, b, "get")
	return Render(c, 200, r.GetResponseParams(), r.GetContentType())
}

// List gets all resources
func List(c *echo.Context, r AzureResource) error {
	groupName := c.Param("group_name")
	resources := make([]map[string]interface{}, 0)
	var parentResources []map[string]interface{}
	var err error

	if groupName != "" {
		// nested route
		parentResources = append(parentResources, map[string]interface{}{"name": groupName})
	} else {
		parentPath := fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, "2015-01-01")
		parentResources, err = GetResources(c, parentPath)
		if err != nil {
			return err
		}
	}

	for _, parent := range parentResources {
		resourcePath := r.GetCollectionPath(parent["name"].(string))
		resp, err := GetResources(c, resourcePath)
		if err != nil {
			return err
		}
		//add href for each resource
		for _, resource := range resp {
			resource["href"] = r.GetHref(parent["name"].(string), resource["name"].(string))
		}
		resources = append(resources, resp...)
	}

	return Render(c, 200, resources, r.GetContentType()+";type=collection")
}

// GetResources makes a call to cloud to get all resources
func GetResources(c *echo.Context, path string) ([]map[string]interface{}, error) {
	client, _ := GetAzureClient(c)
	log.Printf("Get Resources request: %s\n", path)
	resp, err := client.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while requesting resources: %v", err))
	}
	b, _ := ioutil.ReadAll(resp.Body)
	//TODO: add error handling here
	if resp.StatusCode >= 400 {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while requesting resources: %s", string(b)))
	}

	var m map[string][]map[string]interface{}
	var resources []map[string]interface{}
	err = json.Unmarshal(b, &m)
	resources = m["value"]
	if err != nil {
		//try to unmarshal with different interface
		json.Unmarshal(b, &resources)
	}

	return resources, nil
}

// Render sends a JSON resource specific content type response with status code.
func Render(c *echo.Context, code int, resources interface{}, contentType string) error {
	c.Response.Header().Set(echo.ContentType, contentType)
	c.Response.WriteHeader(code)
	return json.NewEncoder(c.Response).Encode(resources)
}
