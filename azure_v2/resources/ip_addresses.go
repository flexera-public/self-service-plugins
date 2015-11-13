package resources

import (
	"encoding/json"
	"fmt"
	"strings"

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
		Name             string `json:"name,omitempty"`
		Location         string `json:"location,omitempty"`
		Group            string `json:"group_name,omitempty"`
		AllocationMethod string `json:"allocation_method,omitempty"` //*Mandatory: Defines whether the IP address is stable or dynamic. Options are Static or Dynamic
		IdleTimeout      int    `json:"timeout,omitempty"`           //Specifies the timeout for the TCP idle connection. The value can be set between 4 and 30 minutes
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
func SetupIPAddressesRoutes(e *echo.Group) {
	e.Get("/ip_addresses", listIPAddresses)
	// e.Get("/ip_addresses/:id", listOneIPAddress)
	// e.Post("/ip_addresses", createIPAddress)
	// e.Delete("/ip_addresses/:id", deleteIPAddress)

	//nested routes
	group := e.Group("/resource_groups/:group_name/ip_addresses")
	group.Get("", listIPAddresses)
	group.Get("/:id", listOneIPAddress)
	group.Post("", createIPAddress)
	group.Delete("/:id", deleteIPAddress)
}

func listIPAddresses(c *echo.Context) error {
	return List(c, new(IPAddress))
}

func listOneIPAddress(c *echo.Context) error {
	ipAddress := IPAddress{
		createParams: ipAddressCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &ipAddress)
}

func createIPAddress(c *echo.Context) error {
	ipAddress := new(IPAddress)
	return Create(c, ipAddress)
}

func deleteIPAddress(c *echo.Context) error {
	ipAddress := IPAddress{
		createParams: ipAddressCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
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
	ip.createParams.Group = c.Param("group_name")
	ip.requestParams.Location = ip.createParams.Location
	ip.requestParams.Properties = map[string]interface{}{
		"publicIPAllocationMethod": ip.createParams.AllocationMethod,
		//Note: dnsSettings could be added if needed
	}
	if ip.createParams.IdleTimeout != 0 {
		ip.requestParams.Properties["idleTimeoutInMinutes"] = ip.createParams.IdleTimeout
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
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, ipAddressPath, config.APIVersion)
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, ipAddressPath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (ip *IPAddress) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &ip.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := ip.GetHref(ip.responseParams.ID)
	if actionName == "create" {
		c.Response().Header().Add("Location", href)
	} else if actionName == "get" {
		ip.responseParams.Href = href
	}
	return nil
}

// GetContentType returns ip address content type
func (ip *IPAddress) GetContentType() string {
	return "vnd.rightscale.ip_address+json"
}

// GetHref returns ip address href
func (ip *IPAddress) GetHref(ipAddressID string) string {
	array := strings.Split(ipAddressID, "/")
	return fmt.Sprintf("/resource_groups/%s/ip_addresses/%s", array[len(array)-5], array[len(array)-1])
}
