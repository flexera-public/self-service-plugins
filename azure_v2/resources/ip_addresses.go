package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	IpAddressPath = "providers/Microsoft.Network/publicIPAddresses"
)

type (
	IpAddressResponseParams struct {
		Id         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	IpAddressRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	IpAddressCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	IpAddress struct {
		CreateParams   IpAddressCreateParams
		RequestParams  IpAddressRequestParams
		ResponseParams IpAddressResponseParams
	}
)

func SetupIpAddressesRoutes(e *echo.Echo) {
	e.Get("/ip_addresses", listIpAddresses)
	e.Get("/ip_addresses/:id", listOneIpAddress)
	e.Post("/ip_addresses", createIpAddress)
	e.Delete("/ip_addresses/:id", deleteIpAddress)

	//nested routes
	//group := e.Group("/resource_groups/:group_name/ip_addresses")
	//group.Get("", listIpAddresses)
	//group.Post("", createIpAddress)
	//group.Delete("/:id", deleteIpAddress)
}

func listIpAddresses(c *echo.Context) error {
	return lib.List(c, new(IpAddress))
}

func listOneIpAddress(c *echo.Context) error {
	params := c.Request.Form
	ip_address := IpAddress{
		CreateParams: IpAddressCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return lib.Get(c, &ip_address)
}

func createIpAddress(c *echo.Context) error {
	ip_address := new(IpAddress)
	return lib.Create(c, ip_address)
}

func deleteIpAddress(c *echo.Context) error {
	params := c.Request.Form
	ip_address := IpAddress{
		CreateParams: IpAddressCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return lib.Delete(c, &ip_address)
}

func (ip *IpAddress) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&ip.CreateParams)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	ip.RequestParams.Location = ip.CreateParams.Location
	ip.RequestParams.Properties = map[string]interface{}{
		"publicIPAllocationMethod": "Dynamic",
		//"dnsSettings":   map[string]interface{}{
		//	"domainNameLabel": postParams.Get("domain_name")
		//}
	}
	return ip.RequestParams, nil
}

func (ip *IpAddress) GetResponseParams() interface{} {
	return ip.ResponseParams
}

func (ip *IpAddress) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, ip.CreateParams.Group, IpAddressPath, ip.CreateParams.Name, config.ApiVersion)
}

func (ip *IpAddress) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, IpAddressPath, config.ApiVersion)
}

func (ip *IpAddress) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &ip.ResponseParams)
	href := ip.GetHref(ip.CreateParams.Group, ip.ResponseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		ip.ResponseParams.Href = href
	}
}

func (ip *IpAddress) GetContentType() string {
	return "vnd.rightscale.ip_address+json"
}

func (ip *IpAddress) GetHref(groupName string, ipAddressName string) string {
	return fmt.Sprintf("/ip_addresses/%s?group_name=%s", groupName, ipAddressName)
}
