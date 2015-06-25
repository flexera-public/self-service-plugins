package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	virtualMachinesPath = "providers/Microsoft.Compute/virtualMachines"
)

type (
	ResponseParams struct {
		Id                string                 `json:"id,omitempty"`
		Type              string                 `json:"type,omitempty"`
		Name              string                 `json:"name"`
		Location          string                 `json:"location"`
		Properties        map[string]interface{} `json:"properties,omitempty"`
		ProvisioningState interface{}            `json:"provisioningState,omitempty"`
		InstanceView      interface{}            `json:"instanceView,omitempty"`
		HardwareProfile   interface{}            `json:"hardwareProfile,omitempty"`
		NetworkProfile    interface{}            `json:"networkProfile,omitempty"`
		StorageProfile    interface{}            `json:"storageProfile,omitempty"`
		Href              string                 `json:"href,omitempty"`
	}

	RequestParams struct {
		Name       string                 `json:"name"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	CreateParams struct {
		Name     string `json:"name,omitempty"`
		Location string `json:"location,omitempty"`
		Size     string `json:"instance_type_uid,omitempty"`
		Group    string `json:"group_name,omitempty"`
	}
	Instance struct {
		CreateParams
		RequestParams
		ResponseParams
	}
)

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
	return List(c, new(Instance))
}
func listOneInstance(c *echo.Context) error {
	params := c.Request.Form
	instance := Instance{
		CreateParams: CreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Get(c, &instance)
}

// https://msdn.microsoft.com/en-us/library/azure/mt163591.aspx
// TODO: check out that provider is already registered - https://msdn.microsoft.com/en-us/library/azure/dn790548.aspx
func createInstance(c *echo.Context) error {
	instance := new(Instance)
	return Create(c, instance)
}

func deleteInstance(c *echo.Context) error {
	params := c.Request.Form
	instance := Instance{
		CreateParams: CreateParams{
			Name:  c.Param("id"),
			Group: params.Get("group_name"),
		},
	}
	return Delete(c, &instance)
}

func (i *Instance) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&i.CreateParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	var networkInterfaces []map[string]interface{}
	i.RequestParams.Name = i.CreateParams.Name
	i.RequestParams.Location = i.CreateParams.Location
	i.RequestParams.Properties = map[string]interface{}{
		"hardwareProfile": map[string]interface{}{"vmSize": i.CreateParams.Size},
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
					"uri": "https://khrvitestgo1.blob.core.windows.net/vhds/cli64174115e3045f57-os-" + i.CreateParams.Name + ".vhd",
				},
				"name": "cli64174115e3045f57-os-" + i.CreateParams.Name,
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
	}
	return i.RequestParams, nil
}

func (i *Instance) GetResponseParams() interface{} {
	return i.ResponseParams
}

func (i *Instance) GetPath() string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, i.CreateParams.Group, virtualMachinesPath, i.CreateParams.Name, "2015-05-01-preview")
}

func (i *Instance) GetCollectionPath(groupName string) string {
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, groupName, virtualMachinesPath, "2015-05-01-preview")
}

func (i *Instance) HandleResponse(c *echo.Context, body []byte, actionName string) {
	json.Unmarshal(body, &i.ResponseParams)
	href := i.GetHref(i.CreateParams.Group, i.ResponseParams.Name)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		i.ResponseParams.Href = href
	}
}

func (i *Instance) GetContentType() string {
	return "vnd.rightscale.instance+json"
}

func (i *Instance) GetHref(groupName string, instanceName string) string {
	return fmt.Sprintf("/instances/%s?group_name=%s", instanceName, groupName)
}
