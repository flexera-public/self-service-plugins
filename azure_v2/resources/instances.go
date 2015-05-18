package resources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/middleware"
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

func SetupInstanceRoutes(e *echo.Echo) {
	e.Get("/instances", listInstances)
	e.Post("/instances", createInstance)
}

func listInstances(c *echo.Context) *echo.HTTPError {
	requestParams := c.Request.Form
	if requestParams.Get("group_name") != "" {
		code, resp := getInstances(c, requestParams.Get("group_name"))
		return c.JSON(code, resp)
	} else {
		code, resp := getResources(c, "")
		var instances []*Instance
		for _, resource_group := range resp {
			_, resp := getInstances(c, resource_group.Name)
			instances = append(instances, resp...)
		}
		// [].to_json => null ... why?
		return c.JSON(code, instances)
	}

}

func getInstances(c *echo.Context, group_name string) (int, []*Instance) {
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, config.ApiVersion)
	log.Printf("Get Instances request: %s\n", path)
	resp, err := client.Get(path)

	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	var m map[string][]*Instance
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	return resp.StatusCode, m["value"]
}

// check out that provider is already registered - https://msdn.microsoft.com/en-us/library/azure/dn790548.aspx
func createInstance(c *echo.Context) *echo.HTTPError {
	postParams := c.Request.Form
	client, _ := middleware.GetAzureClient(c)
	instanceParams := Instance{
		Name:     postParams.Get("name"),
		Location: postParams.Get("location"),
		Properties: map[string]interface{}{
			"hardwareProfile": map[string]interface{}{"vmSize": postParams.Get("instance_type_uid")},
			"storageProfile": map[string]interface{}{
				"osDisk": map[string]interface{}{
					"vhd": map[string]interface{}{
						"uri": "https://khrvi3my1hmm8.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"},
					"name":   "os-" + postParams.Get("name") + "-rs",
					"osType": "Linux"},
				//"destinationVhdsContainer": "http://khrvi.blob.core.windows.net/vhds"}, // hard coded for now...should be used Placement group
			},
			"networkProfile": map[string]interface{}{
				"networkInterfaces": map[string]interface{}{
					"id": "/subscriptions/subscriptionId/resourceGroups/resourceGroupName/providers/Microsoft.Network/NetworkAdapters/Nic1",
				},
			},
		},
	}

	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, postParams.Get("group_name"), virtualMachinesPath, instanceParams.Name, config.ApiVersion)
	log.Printf("Create Instances request with params: %s\n", postParams)
	log.Printf("Create Instances path: %s\n", path)
	data := url.Values{}
	data.Set("name", postParams.Get("name"))
	data.Set("location", "West US")

	by, err := json.Marshal(instanceParams)
	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	log.Printf("READER: %s", reader)
	request, _ := http.NewRequest("PUT", path, reader)
	request.Header.Add("Content-Type", config.MediaType)
	request.Header.Add("Accept", config.MediaType)
	request.Header.Add("User-Agent", config.UserAgent)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Post:", err)
	}
	defer response.Body.Close()
	b, _ := ioutil.ReadAll(response.Body)
	return c.JSON(response.StatusCode, string(b))
}
