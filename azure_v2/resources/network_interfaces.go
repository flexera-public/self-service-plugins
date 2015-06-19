package resources

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	NetworkInterfacePath = "providers/Microsoft.Network/networkInterfaces"
)

type NetworkInterface struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Tags       interface{} `json:"tags,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupNetworkInterfacesRoutes(e *echo.Echo) {
	e.Get("/network_interfaces", listNetworkInterfaces)
	e.Post("/network_interfaces", createNetworkInterface)
	e.Delete("/network_interfaces/:id", deleteNetworkInterface)

	//nested routes
	// group := e.Group("/resource_groups/:group_name/network_interfaces")
	// group.Get("", listNetworkInterfaces)
	//group.Post("", createNetworkInterface)
	//group.Delete("/:id", deleteNetworkInterface)
}

func listNetworkInterfaces(c *echo.Context) error {
	return lib.ListResource(c, NetworkInterfacePath, "network_interfaces")
}

func deleteNetworkInterface(c *echo.Context) error {
	group_name := c.Param("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, NetworkInterfacePath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

func createNetworkInterface(c *echo.Context) error {
	var createParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		SubnetId string `json:"subnet_id,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&createParams)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	var configs []map[string]interface{}
	data := NetworkInterface{
		Location: createParams.Location,
		Properties: map[string]interface{}{
			"ipConfigurations": append(configs, map[string]interface{}{
				"name": createParams.Name + "_ip",
				"properties": map[string]interface{}{
					"subnet": map[string]interface{}{
						"id": createParams.SubnetId,
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
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Group, NetworkInterfacePath, createParams.Name, config.ApiVersion)
	b, err := lib.CreateResource(c, path, data)
	if err != nil {
		return err
	}
	var dat *IpAddress
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}

	return c.JSON(201, dat)
}
