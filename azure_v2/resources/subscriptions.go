package resources

import (
	"fmt"
	"log"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	"github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

const (
	subscriptionsPath = "subscriptions"
)

func SetupSubscriptionRoutes(e *echo.Echo) {
	e.Get("/subscriptions", listSubscriptions)
	e.Get("/subscription", GetSubscription)
}

func listSubscriptions(c *echo.Context) *echo.HTTPError {
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/%s?api-version=%s", config.BaseUrl, subscriptionsPath, config.ApiVersion)
	log.Printf("Get Subscriptions request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return c.JSON(resp.StatusCode, string(body))
}

func GetSubscription(c *echo.Context) *echo.HTTPError {
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/%s/%s?api-version=%s", config.BaseUrl, subscriptionsPath, *config.SubscriptionIdCred, "2015-01-01")
	log.Printf("Get Subscription request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return c.JSON(resp.StatusCode, string(body))
}
