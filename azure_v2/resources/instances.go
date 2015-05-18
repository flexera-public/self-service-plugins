package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

const (
	virtualMachinesPath = "providers/Microsoft.Compute/virtualMachines"
)

type Instance struct {
	provisioningState interface{}
	instanceView      interface{}
	hardwareProfile   interface{}
	networkProfile    interface{}
	storageProfile    interface{}
	id                string `json:"id"`
	name              string `json:"name"`
	Type              string `json:"type"`
	location          string `json:"location"`
}

func SetupInstanceRoutes(e *echo.Echo) {
	e.Get("/instances", listInstances)
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
