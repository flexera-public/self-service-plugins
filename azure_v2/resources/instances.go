package resources

import (
	"log"
	"fmt"
	"io/ioutil"
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
		code, body := getInstances(c, requestParams.Get("group_name"))

		return c.JSON(code, string(body))
	} else {
		code, resp := getResources(c, "")
		var instances []string
		byt := []byte(resp)
		var dat map[string][]ResourceGroup
		if err := json.Unmarshal(byt, &dat); err != nil {
	        log.Fatal("Unmarshaling failed:", err)
	    }
	    log.Printf("Groups: %s\n", len(dat["value"]))
		for _, resource_group := range dat["value"] {
			_, body := getInstances(c, resource_group.Name)
			instances = append(instances, string(body))
		}
		//TODO: concat all responses in one
		return c.JSON(code, instances)
	}
	
}

func getInstances(c *echo.Context, group_name string) (int, string) {
	client, _ := middleware.GetAzureClient(c)
	subscription, _ := middleware.GetCookie(c, middleware.SubscriptionCookieName)
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseUrl, subscription.Value, group_name, virtualMachinesPath, config.ApiVersion)
	log.Printf("Get Instances request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}
