package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	NetworkInterfacePath = "providers/Microsoft.Network/networkInterfaces"
)

type (
	NetworkInterfaceResponseParams struct {
		Id         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Location   string      `json:"location"`
		Tags       interface{} `json:"tags,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	NetworkInterfaceRequestParams struct {
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	NetworkInterfaceCreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		SubnetId string `json:"subnet_id,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	NetworkInterface struct {
		CreateParams   NetworkInterfaceCreateParams
		RequestParams  NetworkInterfaceRequestParams
		ResponseParams NetworkInterfaceResponseParams
	}
)

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
	return lib.List(c, new(NetworkInterface))
}

func listOneNetworkInterface(c *echo.Context) error {
	params := c.Request.Form
	network_interface := NetworkInterface{
		CreateParams: NetworkInterfaceCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return lib.Get(c, &network_interface)
}

func createNetworkInterface(c *echo.Context) error {
	network_interface := new(NetworkInterface)
	return lib.Create(c, network_interface)
}

func deleteNetworkInterface(c *echo.Context) error {
	params := c.Request.Form
	network_interface := NetworkInterface{
		CreateParams: NetworkInterfaceCreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return lib.Delete(c, &network_interface)
}

func (ni *NetworkInterface) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&ni.CreateParams)
	if err != nil {
		return nil, lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	var configs []map[string]interface{}
	ni.RequestParams.Location = ni.CreateParams.Location
	ni.RequestParams.Properties = map[string]interface{}{
		"ipConfigurations": append(configs, map[string]interface{}{
			"name": ni.CreateParams.Name + "_ip",
			"properties": map[string]interface{}{
				"subnet": map[string]interface{}{
					"id": ni.CreateParams.SubnetId,
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

	return ni.RequestParams, nil
}

func (ni *NetworkInterface) GetResponseParams() interface{} {
	return ni.ResponseParams
}

func (ni *NetworkInterface) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, ni.CreateParams.Group, NetworkInterfacePath, ni.CreateParams.Name, config.ApiVersion)
}

func (ni *NetworkInterface) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, NetworkInterfacePath, config.ApiVersion)
}

func (ni *NetworkInterface) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &ni.ResponseParams)
	href := ni.GetHref(ni.CreateParams.Group, ni.ResponseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		ni.ResponseParams.Href = href
	}
}

func (ni *NetworkInterface) GetContentType() string {
	return "vnd.rightscale.network_interface+json"
}

func (ni *NetworkInterface) GetHref(groupName string, interfaceName string) string {
	return fmt.Sprintf("/network_interfaces/%s?group_name=%s", groupName, interfaceName)
}
