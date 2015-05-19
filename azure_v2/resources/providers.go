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

type Provider struct {
	Id                string      `json:"id"`
	Namespace         string      `json:"namespace"`
	RegistrationState string      `json:"registrationState"`
	ResourceTypes     interface{} `json:"resourceTypes"`
	ApplicationID     string      `json:"applicationID,omitempty"`
}

const (
	providerApiVersion = "2015-01-01"
)

func SetupProviderRoutes(e *echo.Echo) {
	e.Get("/providers", listProviders)
	e.Get("/providers/:provider_name", listOneProvider)
	e.Post("/providers/:provider_name/register", registerProvider)
}

func listProviders(c *echo.Context) *echo.HTTPError {
	code, body := getProviders(c, "")
	var dat map[string][]*Provider
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}
	return c.JSON(code, dat["value"])
}

func listOneProvider(c *echo.Context) *echo.HTTPError {
	provider_name := c.Param("provider_name")
	code, body := getProviders(c, provider_name)
	var dat *Provider
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}
	return c.JSON(code, dat)
}

func registerProvider(c *echo.Context) *echo.HTTPError {
	provider_name := c.Param("provider_name")
	_, body := getProviders(c, provider_name)
	var dat *Provider
	if err := json.Unmarshal(body, &dat); err != nil {
		log.Fatal("Unmarshaling failed:", err)
	}
	if dat.RegistrationState == "NotRegistered" {
		log.Printf("Register required: \n")
		client, _ := middleware.GetAzureClient(c)
		path := fmt.Sprintf("%s/subscriptions/%s/providers/%s/register?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, provider_name, providerApiVersion)
		log.Printf("Registering Provider %s: %s\n", provider_name, path)
		resp, err := client.PostForm(path, nil)
		if err != nil {
			log.Fatal("Get:", err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var dat *Provider
		if err := json.Unmarshal(body, &dat); err != nil {
			log.Fatal("Unmarshaling failed:", err)
		}
		return c.JSON(resp.StatusCode, dat)
	}

	return &echo.HTTPError{
		Error: fmt.Errorf("Provider %s already registered.", provider_name),
		Code:  400,
	}
}

func getProviders(c *echo.Context, provider_name string) (int, []byte) {
	client, _ := middleware.GetAzureClient(c)
	path := fmt.Sprintf("%s/subscriptions/%s/providers/%s?api-version=%s", config.BaseUrl, *config.SubscriptionIdCred, provider_name, providerApiVersion)
	log.Printf("Get Providers request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// TODO: handle 400+ statuses
	return resp.StatusCode, body
}
