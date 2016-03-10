package resources

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		Plan       map[string]interface{} `json:"plan,omitempty"`
	}
	createParams struct {
		Name               string                 `json:"name,omitempty"`
		Location           string                 `json:"location,omitempty"`
		Size               string                 `json:"instance_type_uid,omitempty"`
		Group              string                 `json:"group_name,omitempty"`
		NetworkInterfaceID []interface{}          `json:"network_interfaces_ids,omitempty"`
		ImageID            string                 `json:"image_id,omitempty"`
		Plan               map[string]interface{} `json:"image_plan,omitempty"`
		PrivateImageID     string                 `json:"private_image_id,omitempty"`
		StorageAccountID   string                 `json:"storage_account_id,omitempty"`
		HostName           string                 `json:"host_name,omitempty"`
		AdminUserName      string                 `json:"admin_user_name,omitempty"`
		AdminPassword      string                 `json:"admin_password,omitempty"`
		AvailabilitySet    string                 `json:"availability_set,omitempty"`
		Disks              []interface{}          `json:"disks,omitempty"` // [{ "name" : "datadisk1", "diskSizeGB" : "1", "lun" : 0, "vhd":{ "uri" : "http://mystore1.blob.core.windows.net/vhds/dd1.vhd" }, "createOption":"Empty"}]},
		OSDiskName         string                 `json:"os_disk_name,omitempty"`
		UserData           string                 `json:"user_data,omitempty"` // Specifies a base-64 encoded string of custom data. The base-64 encoded string is decoded to a binary array that is saved as a file on the Virtual Machine. The maximum length of the binary array is 65535 bytes.
		WindowsConfig      map[string]interface{} `json:"windows_config,omitempty"`
		LinuxConfig        map[string]interface{} `json:"linux_config,omitempty"`
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
func SetupInstanceRoutes(e *echo.Group) {
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

	group.Put("/:id", updateInstance)
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

	if i.createParams.StorageAccountID == "" {
		return nil, eh.InvalidParamException("storage_account_id")
	}

	if i.createParams.Size == "" {
		return nil, eh.InvalidParamException("instance_type_id")
	}

	i.requestParams.Name = i.createParams.Name
	i.requestParams.Location = i.createParams.Location

	osProfile, err := i.prepareStorageProfile()
	if err != nil {
		return nil, err
	}
	i.requestParams.Properties = map[string]interface{}{
		"hardwareProfile": map[string]interface{}{"vmSize": i.createParams.Size},
		"storageProfile":  osProfile,
		"networkProfile": map[string]interface{}{
			"networkInterfaces": i.createParams.NetworkInterfaceID,
		},
	}

	if i.createParams.PrivateImageID == "" {
		i.requestParams.Properties["osProfile"] = i.prepareOSProfile()
	}

	if i.createParams.AvailabilitySet != "" {
		i.requestParams.Properties["availabilitySet"] = map[string]string{
			"id": i.createParams.AvailabilitySet,
		}
	}

	//Add plan if needed
	if i.createParams.Plan != nil {
		i.requestParams.Plan = i.createParams.Plan
	}
	return i.requestParams, nil
}

func (i *Instance) prepareOSProfile() map[string]interface{} {
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

	osProfile := map[string]interface{}{
		"computerName":  hostName,
		"adminUsername": adminName,
		"adminPassword": adminPassword,
	}

	if i.createParams.WindowsConfig != nil {
		// "windowsConfiguration":{ "provisionVMAgent":true,
		// 		"winRM": { "listeners":[{ "protocol": "https", "certificateUrl": "[parameters('certificateUrl')]"}]},
		//         "additionalUnattendContent":{ "pass":"oobesystem", "component":"Microsoft-Windows-Shell-Setup", "settingName":"FirstLogonCommands|AutoLogon", "content":"<XML unattend content>",},
		//         "enableAutomaticUpdates":true,}
		osProfile["windowsConfiguration"] = i.createParams.WindowsConfig
	}

	if i.createParams.LinuxConfig != nil {
		//"linuxConfiguration":{ "disablePasswordAuthentication":"true|false", "ssh":{ "publicKeys":[{ "path":"Path-Where-To-Place-Public-Key-On-VM", "keyData":"Base64Encoded-public-key-file"}]}}
		osProfile["linuxConfiguration"] = i.createParams.LinuxConfig
	}

	if i.createParams.UserData != "" {
		osProfile["customData"] = base64.StdEncoding.EncodeToString([]byte(i.createParams.UserData))
	}
	return osProfile
}

func (i *Instance) prepareStorageProfile() (map[string]interface{}, error) {
	if i.createParams.ImageID == "" && i.createParams.PrivateImageID == "" {
		return nil, eh.GenericException("One of these two params should be passed: 'image_id' or 'private_image_id'.")
	}
	array := strings.Split(i.createParams.StorageAccountID, "/")
	storageName := array[len(array)-1]
	diskName := i.createParams.OSDiskName

	storageProfile := map[string]interface{}{
		"osDisk": map[string]interface{}{
			"name":         diskName,
			"caching":      "ReadWrite",
			"createOption": "FromImage",
			"vhd": map[string]interface{}{
				"uri": "https://" + storageName + ".blob.core.windows.net/vhds/" + diskName + ".vhd",
			},
		},
	}
	if i.createParams.ImageID != "" {
		array := strings.Split(i.createParams.ImageID, "/")
		if len(array) != 17 {
			return nil, eh.InvalidParamException("image_id")
		}
		publisher := array[8]
		offer := array[12]
		sku := array[14]
		version := array[16]

		storageProfile["imageReference"] = map[string]interface{}{
			"publisher": publisher, //"Canonical",
			"offer":     offer,     //"Ubuntu15.04Snappy",
			"sku":       sku,       //"15.04-Snappy",
			"version":   version,   //"15.04.201505060",
		}
	} else {
		storageProfile["osDisk"].(map[string]interface{})["osType"] = "Linux" // or Windows
		storageProfile["osDisk"].(map[string]interface{})["createOption"] = "attach"
		storageProfile["osDisk"].(map[string]interface{})["vhd"] = map[string]interface{}{
			"uri": i.createParams.PrivateImageID,
		}
	}

	if i.createParams.Disks != nil {
		storageProfile["dataDisks"] = i.createParams.Disks
	}
	return storageProfile, nil
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

	var href string
	if i.action != "getInstanceView" {
		href = i.GetHref(i.responseParams.ID)
	}
	if actionName == "create" {
		c.Response().Header().Add("Location", href)
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
	return fmt.Sprintf("resource_groups/%s/instances/%s", array[len(array)-5], array[len(array)-1])
}

func updateInstance(c *echo.Context) error {
	client, err := GetAzureClient(c)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, c.Param("group_name"), virtualMachinesPath, c.Param("id"), "2015-05-01-preview")
	var bodyParams responseParams
	err = c.Get("bodyDecoder").(*json.Decoder).Decode(&bodyParams)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	by, err := json.Marshal(bodyParams)
	if err != nil {
		eh.GenericException(fmt.Sprintf("Error has occurred while marshaling data: %v", err))
	}

	var reader io.Reader
	reader = bytes.NewBufferString(string(by))
	req, err := http.NewRequest("PUT", path, reader)
	req.Header.Add("Content-Type", config.MediaType)
	req.Header.Add("Accept", config.MediaType)
	req.Header.Add("User-Agent", config.UserAgent)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("Error has occurred while updating instance: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return eh.GenericException(fmt.Sprintf("failed to load response body: %s", err))
	}

	if resp.StatusCode >= 400 {
		return eh.GenericException(fmt.Sprintf("Error has occurred while updating instance: %s", string(body)))
	}

	var response responseParams
	if err := json.Unmarshal(body, &response); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	return Render(c, 200, response, "application/json")
}
