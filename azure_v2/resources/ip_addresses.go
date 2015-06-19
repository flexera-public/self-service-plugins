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
	IpAddressPath = "providers/Microsoft.Network/publicIPAddresses"
)

type IpAddress struct {
	Id         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Location   string      `json:"location"`
	Tags       interface{} `json:"tags,omitempty"`
	Etag       string      `json:"etag,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
}

func SetupIpAddressesRoutes(e *echo.Echo) {
	e.Get("/ip_addresses", listIpAddresses)
	e.Post("/ip_addresses", createIpAddress)
	e.Delete("/ip_addresses/:id", deleteIpAddress)

	//nested routes
	//group := e.Group("/resource_groups/:group_name/ip_addresses")
	//group.Get("", listIpAddresses)
	//group.Post("", createIpAddress)
	//group.Delete("/:id", deleteIpAddress)
}

func listIpAddresses(c *echo.Context) error {
	return lib.ListResource(c, IpAddressPath, "ip_addresses")
}

func deleteIpAddress(c *echo.Context) error {
	postParams := c.Request.Form
	group_name := postParams.Get("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, IpAddressPath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

func createIpAddress(c *echo.Context) error {
	var createParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&createParams)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Group, IpAddressPath, createParams.Name, config.ApiVersion)
	data := IpAddress{
		Location: createParams.Location,
		Properties: map[string]interface{}{
			"publicIPAllocationMethod": "Dynamic",
			//"dnsSettings":   map[string]interface{}{
			//	"domainNameLabel": postParams.Get("domain_name")
			//}
		},
	}

	b, err := lib.CreateResource(c, path, data)
	if err != nil {
		return err
	}
	var dat *IpAddress
	if err := json.Unmarshal(b, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}
	c.Response.Header().Add("Location", "/ip_addresses/"+dat.Name+"?group_name="+createParams.Group)
	return c.NoContent(201)
}
