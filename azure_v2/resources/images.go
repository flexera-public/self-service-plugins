package resources

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

const (
	computePath        = "providers/Microsoft.Compute"
	locationApiVersion = "2016-02-01"
)

// SetupImageRoutes declares routes for Image resource
func SetupImageRoutes(e *echo.Group) {
	e.Get("/locations", listLocations)
	e.Get("/locations/:location/images", listImages)

	//temporal routes
	e.Get("/locations/:location/publishers", listPublishers)
	e.Get("/locations/:location/publishers/:publisher/offers", listOffers)
	e.Get("/locations/:location/publishers/:publisher/offers/:offer/skus", listSkus)
	e.Get("/locations/:location/publishers/:publisher/offers/:offer/skus/:sku/versions", listVersions)
	e.Get("/locations/:location/publishers/:publisher/offers/:offer/skus/:sku/versions/:version", getVersionInfo)
}

func listImages(c *echo.Context) error {
	location := c.Param("location")
	publishers, err := getPublishers(c, location)
	if err != nil {
		return err
	}
	var result []map[string]interface{}
	for _, publisher := range publishers {
		offers, _ := getOffers(c, location, publisher["name"].(string))
		for _, offer := range offers {
			skus, _ := getSkus(c, location, publisher["name"].(string), offer["name"].(string))
			for _, sku := range skus {
				versions, _ := getVersions(c, location, publisher["name"].(string), offer["name"].(string), sku["name"].(string))
				for _, version := range versions {
					version, _ := getVersion(c, location, publisher["name"].(string), offer["name"].(string), sku["name"].(string), version["name"].(string))
					// skip images with invalid version name
					// workaround for versions like - "/Subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/Providers/Microsoft.Compute/Locations/westus/Publishers/Canonical/ArtifactTypes/VMImage/Offers/Ubuntu15.04Snappy/Skus/15.04-Snappy/Versions/15.04.201511272055"
					// "error": {
					//   "code": "InvalidParameter",
					//   "target": "version",
					//   "message": "The value of parameter 'version' is invalid."
					// }
					if version != nil {
						result = append(result, version)
					}
				}
			}
		}
	}

	//TODO: add hrefs or use AzureResource interface
	return c.JSON(200, result)
}

func listLocations(c *echo.Context) error {
	locations, err := getLocations(c)
	if err != nil {
		return err
	}
	return c.JSON(200, locations)
}

func getLocations(c *echo.Context) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/locations?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, locationApiVersion)
	locations, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func listPublishers(c *echo.Context) error {
	location := c.Param("location")
	var locations []map[string]interface{}
	var err error
	if location == "" {
		locations, err = getLocations(c)
		if err != nil {
			return err
		}
	} else {
		locations = []map[string]interface{}{{"name": location}}
	}

	var results []map[string]interface{}
	for _, location := range locations {
		publishers, err := getPublishers(c, location["name"].(string))
		if err != nil {
			return err
		}
		results = append(results, publishers...)
	}
	return c.JSON(200, results)
}

func getPublishers(c *echo.Context, locationName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, computePath, locationName, microsoftComputeApiVersion)
	publishers, err := GetResources(c, path)
	if err != nil {
		fmt.Printf("SKIP FOR %s because of error: %s\n", locationName, err)
		emptyArray := make([]map[string]interface{}, 0)
		return emptyArray, nil
		//return nil, err
	}

	return publishers, nil
}
func listOffers(c *echo.Context) error {
	location := c.Param("location")
	publisher := c.Param("publisher")
	offers, err := getOffers(c, location, publisher)
	if err != nil {
		return err
	}
	return c.JSON(200, offers)
}

func getOffers(c *echo.Context, locationName string, publisherName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, computePath, locationName, publisherName, microsoftComputeApiVersion)
	offers, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

func listSkus(c *echo.Context) error {
	location := c.Param("location")
	publisher := c.Param("publisher")
	offer := c.Param("offer")
	skus, err := getSkus(c, location, publisher, offer)
	if err != nil {
		return err
	}
	return c.JSON(200, skus)
}

func getSkus(c *echo.Context, locationName string, publisherName string, offerName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers/%s/skus?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, computePath, locationName, publisherName, offerName, microsoftComputeApiVersion)
	skus, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return skus, nil
}

func listVersions(c *echo.Context) error {
	location := c.Param("location")
	publisher := c.Param("publisher")
	offer := c.Param("offer")
	sku := c.Param("sku")
	versions, err := getVersions(c, location, publisher, offer, sku)
	if err != nil {
		return err
	}
	return c.JSON(200, versions)
}

func getVersions(c *echo.Context, locationName string, publisherName string, offerName string, skuName string) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers/%s/skus/%s/versions?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, computePath, locationName, publisherName, offerName, skuName, microsoftComputeApiVersion)
	versions, err := GetResources(c, path)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func getVersion(c *echo.Context, locationName string, publisherName string, offerName string, skuName string, versionName string) (map[string]interface{}, error) {
	path := fmt.Sprintf("%s/subscriptions/%s/%s/locations/%s/publishers/%s/artifacttypes/vmimage/offers/%s/skus/%s/versions/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, computePath, locationName, publisherName, offerName, skuName, versionName, microsoftComputeApiVersion)
	body, err := GetResource(c, path)
	if err != nil {
		return nil, err
	}

	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func getVersionInfo(c *echo.Context) error {
	location := c.Param("location")
	publisher := c.Param("publisher")
	offer := c.Param("offer")
	sku := c.Param("sku")
	version := c.Param("version")

	v, err := getVersion(c, location, publisher, offer, sku, version)
	if err != nil {
		return err
	}

	return c.JSON(200, &v)
}
