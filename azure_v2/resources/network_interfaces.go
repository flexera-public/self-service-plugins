package resources

import (
	"encoding/json"
	"fmt"

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
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		SubnetID string `json:"subnet_id,omitempty"`
		Group    string `json:"group_name,omitempty"`
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
	e.Get("/network_interfaces/:id", listOneNetworkInterface)
	e.Post("/network_interfaces", createNetworkInterface)
	e.Delete("/network_interfaces/:id", deleteNetworkInterface)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/network_interfaces")
	// group.Get("", listNetworkInterfaces)
	//group.Post("", createNetworkInterface)
	//group.Delete("/:id", deleteNetworkInterface)
}

func listNetworkInterfaces(c *echo.Context) error {
	return List(c, new(NetworkInterface))
}

func listOneNetworkInterface(c *echo.Context) error {
	params := c.Request.Form
	networkInterface := NetworkInterface{
		createParams: networkInterfaceCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Get(c, &networkInterface)
}

func createNetworkInterface(c *echo.Context) error {
	networkInterface := new(NetworkInterface)
	return Create(c, networkInterface)
}

func deleteNetworkInterface(c *echo.Context) error {
	params := c.Request.Form
	networkInterface := NetworkInterface{
		createParams: networkInterfaceCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
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

	var configs []map[string]interface{}
	ni.requestParams.Location = ni.createParams.Location
	ni.requestParams.Properties = map[string]interface{}{
		"ipConfigurations": append(configs, map[string]interface{}{
			"name": ni.createParams.Name + "_ip",
			"properties": map[string]interface{}{
				"subnet": map[string]interface{}{
					"id": ni.createParams.SubnetID,
				},
				//"privateIPAddress": "10.0.0.8",
				"privateIPAllocationMethod": "Dynamic",
				// "publicIPAddress": map[string]interface{}{
				// 	"id": ""
				// }
			},
		}),
		// "dnsSettings": map[string]interface{}{
		// 	"dnsServers": postParams.Get("dns_servers")
		// }
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
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkInterfacePath, config.APIVersion)
}

// HandleResponse manage raw cloud response
func (ni *NetworkInterface) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &ni.responseParams)
	href := ni.GetHref(ni.createParams.Group, ni.responseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		ni.responseParams.Href = href
	}
}

// GetContentType returns network interface content type
func (ni *NetworkInterface) GetContentType() string {
	return "vnd.rightscale.network_interface+json"
}

// GetHref returns network interface href
func (ni *NetworkInterface) GetHref(groupName string, interfaceName string) string {
	return fmt.Sprintf("/network_interfa/%s?group_name=%s", interfaceName, groupName)
}
