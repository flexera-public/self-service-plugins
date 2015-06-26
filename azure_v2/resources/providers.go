package resources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

// Provider is base struct for Azure Provider resource
type Provider struct {
	ID                string      `json:"id"`
	Namespace         string      `json:"namespace"`
	RegistrationState string      `json:"registrationState"`
	ResourceTypes     interface{} `json:"resourceTypes"`
	ApplicationID     string      `json:"applicationID,omitempty"`
}

const (
	providerAPIVersion = "2015-01-01"
)

// SetupProviderRoutes declares routes for Provider resource
func SetupProviderRoutes(e *echo.Echo) {
	e.Get("/providers", listProviders)
	e.Get("/providers/:provider_name", listOneProvider)
	e.Post("/providers/:provider_name/register", registerProvider)
}

func listProviders(c *echo.Context) error {
	body, err := getProviders(c, "")
	if err != nil {
		return err
	}
	var dat map[string][]*Provider
	if err := json.Unmarshal(body, &dat); err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	return c.JSON(200, dat["value"])
}

func listOneProvider(c *echo.Context) error {
	providerName := c.Param("provider_name")
	body, err := getProviders(c, providerName)
	if err != nil {
		return err
	}
	var dat *Provider
	if err := json.Unmarshal(body, &dat); err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	return c.JSON(200, dat)
}

func registerProvider(c *echo.Context) error {
	providerName := c.Param("provider_name")
	body, err := getProviders(c, providerName)
	if err != nil {
		return err
	}
	var dat *Provider
	if err := json.Unmarshal(body, &dat); err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	if dat.RegistrationState == "NotRegistered" {
		log.Printf("Register required: \n")
		client, err := GetAzureClient(c)
		if err != nil {
			return err
		}
		path := fmt.Sprintf("%s/subscriptions/%s/providers/%s/register?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, providerName, providerAPIVersion)
		log.Printf("Registering Provider %s: %s\n", providerName, path)
		resp, err := client.PostForm(path, nil)
		if err != nil {
			return eh.GenericException(fmt.Sprintf("Error has occurred while registering provider: %v", err))
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
		}

		var dat *Provider
		if err := json.Unmarshal(body, &dat); err != nil {
			return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
		}
		return c.JSON(resp.StatusCode, dat)
	}

	return eh.GenericException(fmt.Sprintf("Provider %s already registered.", providerName))
}

func getProviders(c *echo.Context, providerName string) ([]byte, error) {
	client, err := GetAzureClient(c)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/subscriptions/%s/providers/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, providerName, providerAPIVersion)
	log.Printf("Get Providers request: %s\n", path)
	resp, err := client.Get(path)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while getting provider: %v", err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}
	// TODO: handle 400+ statuses
	return body, nil
}
