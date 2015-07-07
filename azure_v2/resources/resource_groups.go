package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	resourceGroupResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty,omitempty"`
		Location   string      `json:"location,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	resourceGroupRequestParams struct {
		Location string `json:"location"`
	}
	resourceGroupCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
	}
	// ResourceGroup is base struct for Azure Resource Group resource to store input create params,
	// request create params and response params gotten from cloud.
	ResourceGroup struct {
		createParams   resourceGroupCreateParams
		requestParams  resourceGroupRequestParams
		responseParams resourceGroupResponseParams
	}
)

// SetupGroupsRoutes declares routes for resource group resource
func SetupGroupsRoutes(e *echo.Echo) {
	group := e.Group("/resource_groups")
	group.Get("", listResourceGroups)
	group.Get("/:id", listOneResourceGroup)
	group.Post("", createResourceGroup)
	group.Delete("/:id", deleteResourceGroup)
}

func listResourceGroups(c *echo.Context) error {
	return List(c, new(ResourceGroup))
}

func listOneResourceGroup(c *echo.Context) error {
	group := ResourceGroup{
		createParams: resourceGroupCreateParams{
			Name: c.Param("id"),
		},
	}
	return Get(c, &group)
}

func createResourceGroup(c *echo.Context) error {
	group := new(ResourceGroup)
	return Create(c, group)
}

func deleteResourceGroup(c *echo.Context) error {
	group := ResourceGroup{
		createParams: resourceGroupCreateParams{
			Name: c.Param("id"),
		},
	}
	return Delete(c, &group)
}

// GetRequestParams prepares parameters for create resource group request to the cloud
func (rg *ResourceGroup) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&rg.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	rg.requestParams.Location = rg.createParams.Location

	return rg.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (rg *ResourceGroup) GetResponseParams() interface{} {
	return rg.responseParams
}

// GetPath returns full path to the sigle resource group
func (rg *ResourceGroup) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, rg.createParams.Name, "2015-01-01")
}

// GetCollectionPath returns full path to the collection of resource groups
func (rg *ResourceGroup) GetCollectionPath(_ string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, "2015-01-01")
}

// HandleResponse manage raw cloud response
func (rg *ResourceGroup) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &rg.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := rg.GetHref(rg.responseParams.ID)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		rg.responseParams.Href = href
	}
	return nil
}

// GetContentType returns resource group content type
func (rg *ResourceGroup) GetContentType() string {
	return "vnd.rightscale.resource_group+json"
}

// GetHref returns resource group href
func (rg *ResourceGroup) GetHref(groupID string) string {
	array := strings.Split(groupID, "/")
	return fmt.Sprintf("/resource_groups/%s", array[len(array)-1])
}
