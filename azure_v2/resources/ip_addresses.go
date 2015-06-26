package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	ipAddressPath = "providers/Microsoft.Network/publicIPAddresses"
)

type (
	ipAddressResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	ipAddressRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	ipAddressCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	// IPAddress is base struct for Azure Public IP Address resource to store input create params,
	// request create params and response params gotten from cloud.
	IPAddress struct {
		createParams   ipAddressCreateParams
		requestParams  ipAddressRequestParams
		responseParams ipAddressResponseParams
	}
)

// SetupIPAddressesRoutes declares routes for IPAddress resource
func SetupIPAddressesRoutes(e *echo.Echo) {
	e.Get("/ip_addresses", listIPAddresses)
	e.Get("/ip_addresses/:id", listOneIPAddress)
	e.Post("/ip_addresses", createIPAddress)
	e.Delete("/ip_addresses/:id", deleteIPAddress)

	//nested routes
	//group := e.Group("/resource_groups/:group_name/ip_addresses")
	//group.Get("", listIPAddresses)
	//group.Post("", createIPAddress)
	//group.Delete("/:id", deleteIPAddress)
}

func listIPAddresses(c *echo.Context) error {
	return List(c, new(IPAddress))
}

func listOneIPAddress(c *echo.Context) error {
	params := c.Request.Form
	ipAddress := IPAddress{
		createParams: ipAddressCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Get(c, &ipAddress)
}

func createIPAddress(c *echo.Context) error {
	ipAddress := new(IPAddress)
	return Create(c, ipAddress)
}

func deleteIPAddress(c *echo.Context) error {
	params := c.Request.Form
	ipAddress := IPAddress{
		createParams: ipAddressCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Delete(c, &ipAddress)
}

// GetRequestParams prepares parameters for create ip adderss request to the cloud
func (ip *IPAddress) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&ip.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	ip.requestParams.Location = ip.createParams.Location
	ip.requestParams.Properties = map[string]interface{}{
		"publicIPAllocationMethod": "Dynamic",
		//"dnsSettings":   map[string]interface{}{
		//	"domainNameLabel": postParams.Get("domain_name")
		//}
	}
	return ip.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (ip *IPAddress) GetResponseParams() interface{} {
	return ip.responseParams
}

// GetPath returns full path to the sigle ip address
func (ip *IPAddress) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, ip.createParams.Group, ipAddressPath, ip.createParams.Name, config.APIVersion)
}

// GetCollectionPath returns full path to the collection of ip addresses
func (ip *IPAddress) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, ipAddressPath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (ip *IPAddress) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &ip.responseParams)
	href := ip.GetHref(ip.createParams.Group, ip.responseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		ip.responseParams.Href = href
	}
}

// GetContentType returns ip address content type
func (ip *IPAddress) GetContentType() string {
	return "vnd.rightscale.ip_address+json"
}

// GetHref returns ip address href
func (ip *IPAddress) GetHref(groupName string, ipAddressName string) string {
	return fmt.Sprintf("/ip_addresses/%s?group_name=%s", ipAddressName, groupName)
}
