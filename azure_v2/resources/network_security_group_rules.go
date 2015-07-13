package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	eh "github.com/rightscale/self-service-plugins/azure_v2/error_handler"
)

type (
	networkSecurityGroupRuleResponseParams struct {
		ID         string      `json:"id,omitempty"`
		Name       string      `json:"name,omitempty"`
		Etag       string      `json:"etag,omitempty"`
		Properties interface{} `json:"properties,omitempty"`
		Href       string      `json:"href,omitempty"`
	}

	networkSecurityGroupRuleRequestParams struct {
		Properties map[string]interface{} `json:"properties"`
	}
	networkSecurityGroupRuleCreateParams struct {
		Name            string `json:"name,omitempty"`
		Group           string `json:"group_name,omitempty"`
		SecurityGroupID string `json:"security_group_id,omitempty"`
		Type            string

		Description              string `json:"description,omitempty"`                //A description for this rule. Restricted to 140 characters.
		Protocol                 string `json:"protocol,omitempty"`                   // *Mandatory. Network protocol this rule applies to. Can be Tcp, Udp or * to match both.
		SourcePortRange          string `json:"source_port_range,omitempty"`          // *Mandatory. Source Port or Range. Integer or range between 0 and 65535 or * to match any.
		DestinationPortRange     string `json:"destination_port_range,omitempty"`     // *Mandatory. Destination Port or Range. Integer or range between 0 and 65535 or * to match any.
		SourceAddressPrefix      string `json:"source_address_prefix,omitempty"`      // *Mandatory. CIDR or source IP range or * to match any IP. Tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used.
		DestinationAddressPrefix string `json:"destination_address_prefix,omitempty"` // *Mandatory. CIDR or destination IP range or * to match any IP. Tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used.
		Access                   string `json:"access,omitempty"`                     // *Mandatory. Specifies whether network traffic is allowed or denied. Possible values are “Allow” and “Deny”.
		Priority                 int    `json:"priority,omitempty"`                   // *Mandatory. Specifies the priority of the rule. The value can be between 100 and 4096. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.
		Direction                string `json:"direction,omitempty"`                  // *Mandatory. The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values are “Inbound” and “Outbound”.
	}
	// NetworkSecurityGroupRule is base struct for Azure Network Security Group Rule resource to store input create params,
	// request create params and response params gotten from cloud.
	NetworkSecurityGroupRule struct {
		createParams   networkSecurityGroupRuleCreateParams
		requestParams  networkSecurityGroupRuleRequestParams
		responseParams networkSecurityGroupRuleResponseParams
	}
)

// SetupNetworkSecurityGroupRuleRoutes declares routes for NetworkSecurityGroupRule resource
func SetupNetworkSecurityGroupRuleRoutes(e *echo.Echo) {
	e.Get("/network_security_group_rules", listAllNetworkSecurityGroupRules)

	//nested routes
	group := e.Group("/resource_groups/:group_name/network_security_groups/:security_group_name/network_security_group_rules")
	group.Get("", listNetworkSecurityGroupRules)
	group.Get("/:id", listOneNetworkSecurityGroupRule)
	group.Post("", createNetworkSecurityGroupRule)
	group.Delete("/:id", deleteNetworkSecurityGroupRule)

	groupD := e.Group("/resource_groups/:group_name/network_security_groups/:security_group_name/default_network_security_group_rules")
	groupD.Get("", listDefaultNetworkSecurityGroupRules)
	groupD.Get("/:id", listOneDefaultNetworkSecurityGroupRule)
}

func listAllNetworkSecurityGroupRules(c *echo.Context) error {
	rules := make([]map[string]interface{}, 0)
	path := fmt.Sprintf("%s/subscriptions/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, networkSecurityGroupPath, config.APIVersion)
	groups, err := GetResources(c, path)
	if err != nil {
		return err
	}
	for _, group := range groups {
		array := strings.Split(group["id"].(string), "/")
		groupName := array[len(array)-5]
		groupID := group["name"].(string)
		path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/securityRules?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkSecurityGroupPath, groupID, config.APIVersion)
		resp, err := GetResources(c, path)
		if err != nil {
			return err
		}
		for _, rule := range resp {
			rule["href"] = fmt.Sprintf("/resource_groups/%s/network_security_groups/%s/network_security_group_rules/%s", groupName, groupID, rule["name"])
		}
		rules = append(rules, resp...)
	}
	return Render(c, 200, rules, "vnd.rightscale.network_security_group_rule+json;type=collection")
}

func listNetworkSecurityGroupRules(c *echo.Context) error {
	groupName := c.Param("group_name")
	groupID := c.Param("security_group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/securityRules?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkSecurityGroupPath, groupID, config.APIVersion)
	rules, err := GetResources(c, path)
	if err != nil {
		return err
	}
	//add href for each rule
	for _, rule := range rules {
		rule["href"] = fmt.Sprintf("/resource_groups/%s/network_security_groups/%s/network_security_group_rules/%s", groupName, groupID, rule["name"])
	}
	return Render(c, 200, rules, "vnd.rightscale.network_security_group_rule+json;type=collection")
}

func listOneNetworkSecurityGroupRule(c *echo.Context) error {
	rule := NetworkSecurityGroupRule{
		createParams: networkSecurityGroupRuleCreateParams{
			Name:            c.Param("id"),
			Group:           c.Param("group_name"),
			SecurityGroupID: c.Param("security_group_name"),
		},
	}
	return Get(c, &rule)
}

func listOneDefaultNetworkSecurityGroupRule(c *echo.Context) error {
	rule := NetworkSecurityGroupRule{
		createParams: networkSecurityGroupRuleCreateParams{
			Name:            c.Param("id"),
			Group:           c.Param("group_name"),
			SecurityGroupID: c.Param("security_group_name"),
			Type:            "default",
		},
	}
	return Get(c, &rule)
}

func listDefaultNetworkSecurityGroupRules(c *echo.Context) error {
	groupName := c.Param("group_name")
	groupID := c.Param("security_group_name")
	path := fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/defaultSecurityRules?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, groupName, networkSecurityGroupPath, groupID, config.APIVersion)
	rules, err := GetResources(c, path)
	if err != nil {
		return err
	}
	//add href for each rule
	for _, rule := range rules {
		rule["href"] = fmt.Sprintf("/resource_groups/%s/networks/%s/network_security_group_rules/%s", groupName, groupID, rule["name"])
	}
	return Render(c, 200, rules, "vnd.rightscale.network_security_group_rule+json;type=collection")
}

func createNetworkSecurityGroupRule(c *echo.Context) error {
	networkSecurityGroupRule := new(NetworkSecurityGroupRule)
	return Create(c, networkSecurityGroupRule)
}

func deleteNetworkSecurityGroupRule(c *echo.Context) error {
	networkSecurityGroupRule := NetworkSecurityGroupRule{
		createParams: networkSecurityGroupRuleCreateParams{
			Name:            c.Param("id"),
			Group:           c.Param("group_name"),
			SecurityGroupID: c.Param("security_group_name"),
		},
	}
	return Delete(c, &networkSecurityGroupRule)
}

// GetRequestParams prepares parameters for create network security group rule request to the cloud
func (r *NetworkSecurityGroupRule) GetRequestParams(c *echo.Context) (interface{}, error) {
	err := c.Get("bodyDecoder").(*json.Decoder).Decode(&r.createParams)
	if err != nil {
		return nil, eh.GenericException(fmt.Sprintf("Error has occurred while decoding params: %v", err))
	}
	r.createParams.Group = c.Param("group_name")
	r.createParams.SecurityGroupID = c.Param("security_group_name")

	r.requestParams.Properties = map[string]interface{}{
		"description":              r.createParams.Description,
		"protocol":                 r.createParams.Protocol,
		"sourcePortRange":          r.createParams.SourcePortRange,
		"destinationPortRange":     r.createParams.DestinationPortRange,
		"sourceAddressPrefix":      r.createParams.SourceAddressPrefix,
		"destinationAddressPrefix": r.createParams.DestinationAddressPrefix,
		"access":                   r.createParams.Access,
		"priority":                 r.createParams.Priority,
		"direction":                r.createParams.Direction,
	}

	return r.requestParams, nil
}

// GetResponseParams is accessor function for getting access to responseParams struct
func (r *NetworkSecurityGroupRule) GetResponseParams() interface{} {
	return r.responseParams
}

// GetPath returns full path to the sigle network security group rule
func (r *NetworkSecurityGroupRule) GetPath() string {
	resourceName := "securityRules"
	if r.createParams.Type == "default" {
		resourceName = "defaultSecurityRules"
	}
	return fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/%s/%s/%s/%s?api-version=%s", config.BaseURL, *config.SubscriptionIDCred, r.createParams.Group, networkSecurityGroupPath, r.createParams.SecurityGroupID, resourceName, r.createParams.Name, config.APIVersion)
}

// GetCollectionPath is a fake function to support AzureResource by NetworkSecurityGroupRule
func (r *NetworkSecurityGroupRule) GetCollectionPath(groupName string) string { return "" }

// HandleResponse manage raw cloud response
func (r *NetworkSecurityGroupRule) HandleResponse(c *echo.Context, body []byte, actionName string) error {
	if err := json.Unmarshal(body, &r.responseParams); err != nil {
		return eh.GenericException(fmt.Sprintf("got bad response from server: %s", string(body)))
	}
	href := r.GetHref(r.responseParams.ID)
	if actionName == "create" {
		c.Response.Header().Add("Location", href)
	} else if actionName == "get" {
		r.responseParams.Href = href
	}
	return nil
}

// GetContentType returns network security group rule content type
func (r *NetworkSecurityGroupRule) GetContentType() string {
	return "vnd.rightscale.network_security_group_rule+json"
}

// GetHref returns network security group rule href
func (r *NetworkSecurityGroupRule) GetHref(networkSecurityGroupRuleID string) string {
	array := strings.Split(networkSecurityGroupRuleID, "/")
	return fmt.Sprintf("/resource_groups/%s/network_security_groups/%s/network_security_group_rules/%s", array[len(array)-7], array[len(array)-3], array[len(array)-1])
}
