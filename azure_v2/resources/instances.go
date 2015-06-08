package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

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
}

var createParams struct {
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
	Size     string `json:"instance_type_uid,omitempty"`
	Group    string `json:"group_name,omitempty"`
}

func SetupInstanceRoutes(e *echo.Echo) {
	//get all instances from all groups
	e.Get("/instances", listInstances)
	e.Get("/instances/:id", listOneInstance)
	e.Post("/instances", createInstance)

	//nested routes
	group := e.Group("/resource_groups/:group_name/instances")
	group.Get("", listInstances)
	//group.Post("", createInstance)
	group.Delete("/:id", deleteInstance)
}

func listInstances(c *echo.Context) error {
	return lib.ListResource(c, virtualMachinesPath)
}
func listOneInstance(c *echo.Context) error {
	params := c.Request.Form
	group_name := params.Get("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, c.Param("id"), config.ApiVersion)
	resource, err := lib.GetResource(c, path)
	if err != nil {
		return err
	}
	resource["href"] = "/azure_plugin/instances/" + c.Param("id") + "?group_name=" + group_name
	c.Response.Header().Add("Content-Type", "vnd.rightscale.instance+json")
	return c.JSON(200, resource)
}

func deleteInstance(c *echo.Context) error {
	group_name := c.Param("group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, c.Param("id"), config.ApiVersion)
	return lib.DeleteResource(c, path)
}

// check out that provider is already registered - https://msdn.microsoft.com/en-us/library/azure/dn790548.aspx
func createInstance(c *echo.Context) error {
	var str string
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&str) //decode body as a string due to bug in SS
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	json.Unmarshal([]byte(str), &createParams)
	client, _ := lib.GetAzureClient(c)
	var networkInterfaces []map[string]interface{}
	instanceParams := Instance{
		Name:     createParams.Name,
		Location: createParams.Location,
		Properties: map[string]interface{}{
			"hardwareProfile": map[string]interface{}{"vmSize": createParams.Size},
			"storageProfile": map[string]interface{}{
				"osDisk": map[string]interface{}{
					"vhd": map[string]interface{}{
						"uri": "https://khrvitestgo.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"},
					"name":   "os-" + createParams.Name + "-rs",
					"osType": "Linux"},
				// "sourceImage": map[string]interface{}{
				// 	"id": "/2d2b2267-ff0a-46d3-9912-8577acb18a0a/services/images/7bb63e06fb004b2597e854325d2fe7b9__Test-Windows-Server-2012-Datacenter-201401.01-en.us-127GB.vhd",
				// },
				// "destinationVhdsContainer": "http://khrvitestgo.blob.core.windows.net/vhds", // hard coded for now...should be used Placement group
			},
			"networkProfile": map[string]interface{}{
				"networkInterfaces": append(networkInterfaces, map[string]interface{}{
					"id": "/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni",
				}),
			},
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, createParams.Group, virtualMachinesPath, instanceParams.Name, config.ApiVersion)
	log.Printf("Create Instances request with params: %#v\n", createParams)
	log.Printf("Create Instances path: %s\n", path)

	by, err := json.Marshal(instanceParams)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	request, err := http.NewRequest("PUT", path, reader)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while creating instance: %v", err))
	}
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		return lib.GenericException(fmt.Sprintf("Error has occurred while creating instance: %v", err))
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return lib.GenericException(fmt.Sprintf("Error has occurred while creating instance: %s", string(b)))
	}

	var m *Instance
	json.Unmarshal(b, &m)
	c.Response.Header().Add("Location", "/azure_plugin/instances/"+m.Name+"?group_name="+createParams.Group)
	fmt.Printf("HEADER: %s", "/azure_plugin/instances/"+m.Name+"?group_name="+createParams.Group)
	return c.JSON(response.StatusCode, "")
}
