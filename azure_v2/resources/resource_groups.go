package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	ResourceGroupResponseParams struct {
		Id         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty,omitempty"`
		Location   string      `json:"location,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	ResourceGroupRequestParams struct {
		Location string `json:"location"`
	}
	ResourceGroupCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
	}
	ResourceGroup struct {
		CreateParams   ResourceGroupCreateParams
		RequestParams  ResourceGroupRequestParams
		ResponseParams ResourceGroupResponseParams
	}
)

func SetupGroupsRoutes(e *echo.Echo) {
	e.Get("/resource_groups", listResourceGroups)
	e.Get("/resource_groups/:id", listOneResourceGroup)
	e.Post("/resource_groups", createResourceGroup)
	e.Delete("/resource_groups/:id", deleteResourceGroup)
}

func listResourceGroups(c *echo.Context) error {
	return List(c, new(ResourceGroup))
}

func listOneResourceGroup(c *echo.Context) error {
	group := ResourceGroup{
		CreateParams: ResourceGroupCreateParams{
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
		CreateParams: ResourceGroupCreateParams{
			Name: c.Param("id"),
		},
	}
	return Delete(c, &group)
}

func (rg *ResourceGroup) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&rg.CreateParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	rg.RequestParams.Location = rg.CreateParams.Location

	return rg.RequestParams, nil
}

func (rg *ResourceGroup) GetResponseParams() interface{} {
	return rg.ResponseParams
}

func (rg *ResourceGroup) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, rg.CreateParams.Name, "2015-01-01")
}

func (rg *ResourceGroup) GetCollectionPath(_ string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, "2015-01-01")
}

func (rg *ResourceGroup) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &rg.ResponseParams)
	href := rg.GetHref(rg.ResponseParams.Name, "")
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		rg.ResponseParams.Href = href
	}
}

func (rg *ResourceGroup) GetContentType() string {
	return "vnd.rightscale.resource_group+json"
}

func (rg *ResourceGroup) GetHref(groupName string, _ string) string {
	return fmt.Sprintf("/resource_groups/%s", groupName)
}
