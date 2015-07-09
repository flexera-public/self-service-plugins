package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

const (
	virtualMachinesPath  = "providers/Microsoft.Compute/virtualMachines"
	defaultAdminUserName = "rsadministrator"
	defaultAdminPassword = "Pass1234@"
)

type (
	responseParams struct {
		ID                string                 `json:"id,omitempty"`
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

	instanceViewResponseParams struct {
		PlatformUpdateDomain int
		PlatformFaultDomain  int
		VMAgent              map[string]interface{}
		Disks                []interface{}
		Statuses             []interface{}
	}

	requestParams struct {
		Name       string                 `json:"name"`
		Location   string                 `json:"location"`
		Properties map[string]interface{} `json:"properties,omitempty"`
	}
	createParams struct {
		Name               string `json:"name,omitempty"`
		Location           string `json:"location,omitempty"`
		Size               string `json:"instance_type_uid,omitempty"`
		Group              string `json:"group_name,omitempty"`
		NetworkInterfaceID string `json:"network_interface_id,omitempty"`
		ImageID            string `json:"image_id,omitempty"`
		StorageAccountID   string `json:"storage_account_id,omitempty"`
		HostName           string `json:"host_name,omitempty"`
		AdminUserName      string `json:"admin_user_name,omitempty"`
		AdminPassword      string `json:"admin_password,omitempty"`
		AvailabilitySet    string `json:"availability_set,omitempty"`
	}

	// Instance is base struct for Azure VM resource to store input create params,
	// request create params and response params gotten from cloud.
	Instance struct {
		action string
		createParams
		requestParams
		responseParams
		instanceViewResponseParams
	}
)

// SetupInstanceRoutes declares routes for Instance resource
func SetupInstanceRoutes(e *echo.Echo) {
	//get all instances from all groups
	e.Get("/instances", listInstances)
	// e.Get("/instances/:id", listOneInstance)
	// e.Post("/instances", createInstance)
	// e.Delete("/instances/:id", deleteInstance)

	//nested routes
	group := e.Group("/resource_groups/:group_name/instances")
	group.Get("", listInstances)
	group.Get("/:id", listOneInstance)
	group.Get("/:id/instance_view", listOneInstanceView)
	group.Post("", createInstance)
	group.Delete("/:id", deleteInstance)
}

func listInstances(c *echo.Context) error {
	return List(c, new(Instance))
}
func listOneInstance(c *echo.Context) error {
	instance := Instance{
		createParams: createParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Get(c, &instance)
}

func listOneInstanceView(c *echo.Context) error {
	instance := Instance{
		action: "getInstanceView",
		createParams: createParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
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
	instance := Instance{
		createParams: createParams{
			Name:  c.Param("id"),
			Group: c.Param("group_name"),
		},
	}
	return Delete(c, &instance)
}

// GetRequestParams prepares parameters for create instance request to the cloud
func (i *Instance) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&i.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	i.createParams.Group = c.Param("group_name")

	//TODO: make a func for validating createParams and return all errors at once
	if i.createParams.Name == "" {
		return nil, eh.InvalidParamException("name")
	}

	if i.createParams.Location == "" {
		return nil, eh.InvalidParamException("location")
	}

	if i.createParams.ImageID == "" {
		return nil, eh.InvalidParamException("image_id")
	}

	array := strings.Split(i.createParams.ImageID, "/")
	if len(array) != 17 {
		return nil, eh.InvalidParamException("image_id")
	}

	publisher := array[8]
	offer := array[12]
	sku := array[14]
	version := array[16]

	if i.createParams.StorageAccountID == "" {
		return nil, eh.InvalidParamException("storage_account_id")
	}

	if i.createParams.Size == "" {
		return nil, eh.InvalidParamException("instance_type_id")
	}

	array = strings.Split(i.createParams.StorageAccountID, "/")
	storageName := array[len(array)-1]

	hostName := i.createParams.HostName
	if hostName == "" {
		hostName = i.createParams.Name
	}

	adminName := i.createParams.AdminUserName
	if adminName == "" {
		adminName = defaultAdminUserName
	}
	adminPassword := i.createParams.AdminPassword
	if adminPassword == "" {
		adminPassword = defaultAdminPassword
	}

	var networkInterfaces []map[string]interface{}
	i.requestParams.Name = i.createParams.Name
	i.requestParams.Location = i.createParams.Location
	i.requestParams.Properties = map[string]interface{}{
		"hardwareProfile": map[string]interface{}{"vmSize": i.createParams.Size},
		//"storageProfile":{"imageReference":{"publisher":"CoreOS","offer":"CoreOS","sku":"Alpha","version":"660.0.0"},"osDisk":{"name":"cli64174115e3045f57-os-1434041634239","vhd":{"uri":"https://cli64174115e3045f5714340.blob.core.windows.net/vhds/cli64174115e3045f57-os-1434041634239.vhd"},"caching":"ReadWrite","createOption":"FromImage"}}
		"storageProfile": map[string]interface{}{
			"imageReference": map[string]interface{}{
				"publisher": publisher, //"Canonical",
				"offer":     offer,     //"Ubuntu15.04Snappy",
				"sku":       sku,       //"15.04-Snappy",
				"version":   version,   //"15.04.201505060",
			},

			"osDisk": map[string]interface{}{
				"caching":      "ReadWrite",
				"createOption": "FromImage",
				"vhd": map[string]interface{}{
					"uri": "https://" + storageName + ".blob.core.windows.net/vhds/os-" + i.createParams.Name + "-rs.vhd",
				},
				"name": "os-" + i.createParams.Name + "-rs",
				//"osType": "Windows",
			},
		},
		"osProfile": map[string]interface{}{
			"computerName":  hostName,
			"adminUsername": adminName,
			"adminPassword": adminPassword,
			//"linuxConfiguration":{"disablePasswordAuthentication":false}
		},
		"networkProfile": map[string]interface{}{
			"networkInterfaces": append(networkInterfaces, map[string]interface{}{
				"id": i.createParams.NetworkInterfaceID,
			}),
		},
	}

	if i.createParams.AvailabilitySet != "" {
		i.requestParams.Properties["availabilitySet"] = map[string]string{
			"id": i.createParams.AvailabilitySet,
		}
	}
	return i.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (i *Instance) GetResponseParams() interface{} {
	if i.action == "getInstanceView" {
		return i.instanceViewResponseParams
	}
	return i.responseParams
}

// GetPath returns full path to the sigle instance
func (i *Instance) GetPath() string {
	if i.action == "getInstanceView" {
		return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/InstanceView?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, i.createParams.Group, virtualMachinesPath, i.createParams.Name, "2015-05-01-preview")
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, i.createParams.Group, virtualMachinesPath, i.createParams.Name, "2015-05-01-preview")
}

// GetCollectionPath returns full path to the collection of instances
func (i *Instance) GetCollectionPath(groupName string) string {
	if groupName == "" {
		return fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, virtualMachinesPath, "2015-05-01-preview")
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, virtualMachinesPath, "2015-05-01-preview")
}

// HandleResponse manage raw cloud response
func (i *Instance) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	var err error
	if i.action == "getInstanceView" {
		err = json.Unmarshal(body, &i.instanceViewResponseParams)
	} else {
		err = json.Unmarshal(body, &i.responseParams)
	}
	if err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}

	href := i.GetHref(i.responseParams.ID)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		i.responseParams.Href = href
	}
	return nil
}

// GetContentType returns instance content type
func (i *Instance) GetContentType() string {
	return "vnd.rightscale.instance+json"
}

// GetHref returns instance href
func (i *Instance) GetHref(instanceID string) string {
	array := strings.Split(instanceID, "/")
	return fmt.Sprintf("/resource_groups/%s/instances/%s", array[len(array)-5], array[len(array)-1])
}
