package resources

import (
	"fmt"
	"io/ioutil"
	"log"
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

const (
	virtualMachinesPath = "providers/Microsoft.Compute/virtualMachines"
)

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
		var instances []interface{}
		byt := []byte(resp)
		var dat map[string][]*ResourceGroup
		if err := json.Unmarshal(byt, &dat); err != nil {
			log.Fatal("Unmarshaling failed:", err)
		}
		log.Printf("Groups: %s\n", len(dat["value"]))
		for _, resource_group := range dat["value"] {
			_, resp := getInstances(c, resource_group.Name)
			instances = append(instances, resp)
		}
		//TODO: concat all responses in one
		return c.JSON(code, instances)
	}

}

func getInstances(c *echo.Context, group_name string) (int, interface{}) {
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, group_name, virtualMachinesPath, config.ApiVersion)
	log.Printf("Get Instances request: %s\n", path)
	resp, err := client.Get(path)

	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	var m interface{}
	b, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(b, &m)

	return resp.StatusCode, m
}
