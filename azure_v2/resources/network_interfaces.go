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
	networkInterfacePath = "providers/Microsoft.Network/networkInterfaces"
)

type (
	networkInterfaceResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	networkInterfaceRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	networkInterfaceCreateParams struct {
		Name                   string   `json:"name,omitempty"`
		Location               string   `json:"location,omitempty"`
		SubnetID               string   `json:"subnet_id,omitempty"`
		Group                  string   `json:"group_name,omitempty"`
		DNSServers             []string `json:"dns_servers,omitempty"`
		NetworkSecurityGroupID string   `json:"network_security_group_id,omitempty"`
		PrivateIPAddress       string   `json:"private_ip_address,omitempty"`   // Static IP Address
		PublicIPAddressID      string   `json:"public_ip_address_id,omitempty"` // Reference to a Public IP Address to associate with this NIC
	}
	// NetworkInterface is base struct for Azure Network Interface resource to store input create params,
	// request create params and response params gotten from cloud.
	NetworkInterface struct {
		createParams   networkInterfaceCreateParams
		requestParams  networkInterfaceRequestParams
		responseParams networkInterfaceResponseParams
	}
)

// SetupNetworkInterfacesRoutes declares routes for IPAddress resource
func SetupNetworkInterfacesRoutes(e *echo.Echo) {
	e.Get("/network_interfaces", listNetworkInterfaces)
	// e.Get("/network_interfaces/:id", listOneNetworkInterface)
	// e.Post("/network_interfaces", createNetworkInterface)
	// e.Delete("/network_interfaces/:id", deleteNetworkInterface)

	//nested routes
	group := e.Group("/resource_groups/:group_name/network_interfaces")
	group.Get("", listNetworkInterfaces)
	group.Get("/:id", listOneNetworkInterface)
	group.Post("", createNetworkInterface)
	group.Delete("/:id", deleteNetworkInterface)
}

func listNetworkInterfaces(c *echo.Context) error {
	return List(c, new(NetworkInterface))
}

func listOneNetworkInterface(c *echo.Context) error {
	networkInterface := NetworkInterface{
		createParams: networkInterfaceCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &networkInterface)
}

func createNetworkInterface(c *echo.Context) error {
	networkInterface := new(NetworkInterface)
	return Create(c, networkInterface)
}

func deleteNetworkInterface(c *echo.Context) error {
	networkInterface := NetworkInterface{
		createParams: networkInterfaceCreateParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Delete(c, &networkInterface)
}

// GetRequestParams prepares parameters for create network interface request to the cloud
func (ni *NetworkInterface) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&ni.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	ni.createParams.Group = c.Param("group_name")

	ni.requestParams.Location = ni.createParams.Location
	configProperties := map[string]interface{}{
		"subnet": map[string]interface{}{
			"id": ni.createParams.SubnetID,
		},
		"privateIPAllocationMethod": "Dynamic",
	}

	if ni.createParams.PrivateIPAddress != "" {
		configProperties["privateIPAddress"] = ni.createParams.PrivateIPAddress
		configProperties["privateIPAllocationMethod"] = "Static"
	}

	if ni.createParams.PublicIPAddressID != "" {
		configProperties["publicIPAddress"] = map[string]interface{}{
			"id": ni.createParams.PublicIPAddressID,
		}
	}

	ni.requestParams.Properties = map[string]interface{}{
		"ipConfigurations": []map[string]interface{}{
			{"name": ni.createParams.Name + "_ip",
				"properties": configProperties},
		},
	}

	if ni.createParams.NetworkSecurityGroupID != "" {
		ni.requestParams.Properties["networkSecurityGroup"] = map[string]interface{}{
			"id": ni.createParams.NetworkSecurityGroupID,
		}
	}

	if ni.createParams.DNSServers != nil {
		ni.requestParams.Properties["dnsSettings"] = map[string]interface{}{
			"dnsServers": ni.createParams.DNSServers,
		}
	}

	return ni.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (ni *NetworkInterface) GetResponseParams() interface{} {
	return ni.responseParams
}

// GetPath returns full path to the sigle network interface
func (ni *NetworkInterface) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, ni.createParams.Group, networkInterfacePath, ni.createParams.Name, config.APIVersion)
}

// GetCollectionPath returns full path to the collection of network interfaces
func (ni *NetworkInterface) GetCollectionPath(groupName string) string {
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, networkInterfacePath, config.APIVersion)
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkInterfacePath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (ni *NetworkInterface) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &ni.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := ni.GetHref(ni.responseParams.ID)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		ni.responseParams.Href = href
	}
	return nil
}

// GetContentType returns network interface content type
func (ni *NetworkInterface) GetContentType() string {
	return "vnd.rightscale.network_interface+json"
}

// GetHref returns network interface href
func (ni *NetworkInterface) GetHref(networkInterfaceID string) string {
	array := strings.Split(networkInterfaceID, "/")
	return fmt.Sprintf("/resource_groups/%s/network_interfaces/%s", array[len(array)-5], array[len(array)-1])
}
