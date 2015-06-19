package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/lib"
)

const (
	virtualMachinesPath = "providers/Microsoft.Compute/virtualMachines"
)

type Instance struct {
	ProvisioningState interface{}            `json:"provisioningState,omitempty"`
	InstanceView      interface{}            `json:"instanceView,omitempty"`
	HardwareProfile   interface{}            `json:"hardwareProfile,omitempty"`
	NetworkProfile    interface{}            `json:"networkProfile,omitempty"`
	StorageProfile    interface{}            `json:"storageProfile,omitempty"`
	Id                string                 `json:"id,omitempty"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type,omitempty"`
	Location          string                 `json:"location"`
	Properties        map[string]interface{} `json:"properties,omitempty"` // used for create instance
	// Plan              map[string]interface{} `json:"plan,omitempty"`       // used for create instance
}

func SetupInstanceRoutes(e *echo.Echo) {
	//get all instances from all groups
	e.Get("/instances", listInstances)
	e.Get("/instances/:id", listOneInstance)
	e.Post("/instances", createInstance)
	e.Delete("/instances/:id", deleteInstance)

	//nested routes
	group := e.Group("/resource_groups/:group_name/instances")
	group.Get("", listInstances)
	//group.Post("", createInstance)
	//group.Delete("/:id", deleteInstance)
}

func listInstances(c *echo.Context) error {
	return lib.ListResource(c, virtualMachinesPath, "instances")
}
func listOneInstance(c *echo.Context) error {
	params := c.Request.Form
	group_name := params.Get("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, c.Param("id"), config.ApiVersion)
	resource, err := lib.GetResource(c, path, "/instances/%s?group_name="+group_name)
	if err != nil {
		return err
	}

	return lib.Render(c, 200, resource, "vnd.rightscale.instance+json")
}

func deleteInstance(c *echo.Context) error {
	params := c.Request.Form
	group_name := params.Get("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

// https://msdn.microsoft.com/en-us/library/azure/mt163591.aspx
// TODO: check out that provider is already registered - https://msdn.microsoft.com/en-us/library/azure/dn790548.aspx
func createInstance(c *echo.Context) error {
	var createParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Size     string `json:"instance_type_uid,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&createParams)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}

	var networkInterfaces []map[string]interface{}
	instanceParams := Instance{
		Name:     createParams.Name,
		Location: createParams.Location,
		// Plan is required for images from marketplace
		// Plan: map[string]interface{}{
		// 	"name":      "Ubuntu15.04Snappy",
		// 	"publisher": "Canonical",
		// 	//"product":   "imageProduct",
		// },
		Properties: map[string]interface{}{
			"hardwareProfile": map[string]interface{}{"vmSize": createParams.Size},
			//"storageProfile":{"imageReference":{"publisher":"CoreOS","offer":"CoreOS","sku":"Alpha","version":"660.0.0"},"osDisk":{"name":"cli64174115e3045f57-os-1434041634239","vhd":{"uri":"https://cli64174115e3045f5714340.blob.core.windows.net/vhds/cli64174115e3045f57-os-1434041634239.vhd"},"caching":"ReadWrite","createOption":"FromImage"}}
			"storageProfile": map[string]interface{}{
				"imageReference": map[string]interface{}{
					"publisher": "Canonical",         //publisher,
					"offer":     "Ubuntu15.04Snappy", //offer,
					"sku":       "15.04-Snappy",      //sku,
					"version":   "15.04.201505060",   //version,
				},

				"osDisk": map[string]interface{}{
					"caching":      "ReadWrite",
					"createOption": "FromImage",
					"vhd": map[string]interface{}{
						"uri": "https://khrvitestgo1.blob.core.windows.net/vhds/cli64174115e3045f57-os-" + createParams.Name + ".vhd",
					},
					"name": "cli64174115e3045f57-os-" + createParams.Name,
					//"osType": "Windows",
				},
			},
			"osProfile": map[string]interface{}{
				"computerName":  "khrvi",
				"adminUsername": "azureuser",
				"adminPassword": "Pass1234@",
				//"linuxConfiguration":{"disablePasswordAuthentication":false}
			},
			"networkProfile": map[string]interface{}{
				"networkInterfaces": append(networkInterfaces, map[string]interface{}{
					"id": "/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni",
				}),
			},
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Group, virtualMachinesPath, instanceParams.Name, "2015-05-01-preview")
	b, err := lib.CreateResource(c, path, instanceParams)
	if err != nil {
		return err
	}
	var m *Instance
	json.Unmarshal(b, &m)
	c.Response.Header().Add("Location", "/instances/"+m.Name+"?group_name="+createParams.Group)
	return c.NoContent(201)
}
