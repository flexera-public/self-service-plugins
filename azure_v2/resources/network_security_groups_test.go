package resources

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
)

const (
	listNSGsResponse   = `{"value": [{"etag":"W/\"3d45798e-89d0-4c9b-a65d-eecf2b44e6f6\"","href":"/resource_groups/Group-1/network_security_groups/khrvi1","id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1","location":"westus","name":"khrvi1","properties":{"defaultSecurityRules":[{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowVnetInBound","name":"AllowVnetInBound","properties":{"access":"Allow","description":"Allow inbound traffic from all VMs in VNET","destinationAddressPrefix":"VirtualNetwork","destinationPortRange":"*","direction":"Inbound","priority":65000,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"VirtualNetwork","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowAzureLoadBalancerInBound","name":"AllowAzureLoadBalancerInBound","properties":{"access":"Allow","description":"Allow inbound traffic from azure load balancer","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Inbound","priority":65001,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"AzureLoadBalancer","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/DenyAllInBound","name":"DenyAllInBound","properties":{"access":"Deny","description":"Deny all inbound traffic","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Inbound","priority":65500,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowVnetOutBound","name":"AllowVnetOutBound","properties":{"access":"Allow","description":"Allow outbound traffic from all VMs to all VMs in VNET","destinationAddressPrefix":"VirtualNetwork","destinationPortRange":"*","direction":"Outbound","priority":65000,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"VirtualNetwork","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowInternetOutBound","name":"AllowInternetOutBound","properties":{"access":"Allow","description":"Allow outbound traffic from all VMs to Internet","destinationAddressPrefix":"Internet","destinationPortRange":"*","direction":"Outbound","priority":65001,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/DenyAllOutBound","name":"DenyAllOutBound","properties":{"access":"Deny","description":"Deny all outbound traffic","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Outbound","priority":65500,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}}],"provisioningState":"Succeeded"}},{"etag":"W/\"f16517fd-a1f6-449d-a77c-06ef63318423\"","href":"/resource_groups/Group-1/network_security_groups/khrvi2","id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2","location":"westus","name":"khrvi2","properties":{"defaultSecurityRules":[{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/defaultSecurityRules/AllowVnetInBound","name":"AllowVnetInBound","properties":{"access":"Allow","description":"Allow inbound traffic from all VMs in VNET","destinationAddressPrefix":"VirtualNetwork","destinationPortRange":"*","direction":"Inbound","priority":65000,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"VirtualNetwork","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/defaultSecurityRules/AllowAzureLoadBalancerInBound","name":"AllowAzureLoadBalancerInBound","properties":{"access":"Allow","description":"Allow inbound traffic from azure load balancer","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Inbound","priority":65001,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"AzureLoadBalancer","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/defaultSecurityRules/DenyAllInBound","name":"DenyAllInBound","properties":{"access":"Deny","description":"Deny all inbound traffic","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Inbound","priority":65500,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/defaultSecurityRules/AllowVnetOutBound","name":"AllowVnetOutBound","properties":{"access":"Allow","description":"Allow outbound traffic from all VMs to all VMs in VNET","destinationAddressPrefix":"VirtualNetwork","destinationPortRange":"*","direction":"Outbound","priority":65000,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"VirtualNetwork","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/defaultSecurityRules/AllowInternetOutBound","name":"AllowInternetOutBound","properties":{"access":"Allow","description":"Allow outbound traffic from all VMs to Internet","destinationAddressPrefix":"Internet","destinationPortRange":"*","direction":"Outbound","priority":65001,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/defaultSecurityRules/DenyAllOutBound","name":"DenyAllOutBound","properties":{"access":"Deny","description":"Deny all outbound traffic","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Outbound","priority":65500,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}}],"provisioningState":"Succeeded","securityRules":[{"etag":"W/\"f16517fd-a1f6-449d-a77c-06ef63318423\"","id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/securityRules/khrvi2_1","name":"khrvi2_1","properties":{"access":"Allow","description":"test","destinationAddressPrefix":"*","destinationPortRange":"800","direction":"Inbound","priority":100,"protocol":"Tcp","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"800"}}]}}]}`
	listOneNSGResponse = `{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1","name":"khrvi1","location":"westus","etag":"W/\"3d45798e-89d0-4c9b-a65d-eecf2b44e6f6\"","properties":{"defaultSecurityRules":[{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowVnetInBound","name":"AllowVnetInBound","properties":{"access":"Allow","description":"Allow inbound traffic from all VMs in VNET","destinationAddressPrefix":"VirtualNetwork","destinationPortRange":"*","direction":"Inbound","priority":65000,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"VirtualNetwork","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowAzureLoadBalancerInBound","name":"AllowAzureLoadBalancerInBound","properties":{"access":"Allow","description":"Allow inbound traffic from azure load balancer","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Inbound","priority":65001,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"AzureLoadBalancer","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/DenyAllInBound","name":"DenyAllInBound","properties":{"access":"Deny","description":"Deny all inbound traffic","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Inbound","priority":65500,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowVnetOutBound","name":"AllowVnetOutBound","properties":{"access":"Allow","description":"Allow outbound traffic from all VMs to all VMs in VNET","destinationAddressPrefix":"VirtualNetwork","destinationPortRange":"*","direction":"Outbound","priority":65000,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"VirtualNetwork","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/AllowInternetOutBound","name":"AllowInternetOutBound","properties":{"access":"Allow","description":"Allow outbound traffic from all VMs to Internet","destinationAddressPrefix":"Internet","destinationPortRange":"*","direction":"Outbound","priority":65001,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}},{"id":"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1/defaultSecurityRules/DenyAllOutBound","name":"DenyAllOutBound","properties":{"access":"Deny","description":"Deny all outbound traffic","destinationAddressPrefix":"*","destinationPortRange":"*","direction":"Outbound","priority":65500,"protocol":"*","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"*"}}],"provisioningState":"Succeeded"},"href":"/resource_groups/Group-1/network_security_groups/khrvi1"}`
)

var _ = Describe("network_security_groups", func() {

	var do *ghttp.Server
	var client *AzureClient
	var response *Response
	var err error

	BeforeEach(func() {
		do = ghttp.NewServer()
		config.BaseURL = do.URL()
		client = NewAzureClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("listing", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath),
					ghttp.RespondWith(http.StatusOK, listNSGsResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/network_security_groups")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_security_group+json;type=collection"))
		})

		It("lists all network_security_groups inside one resource group", func() {
			networkSecurityGroups := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNSGsResponse), &networkSecurityGroups)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkSecurityGroups["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			subscriptionID := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/"+networkSecurityGroupPath),
					ghttp.RespondWith(http.StatusOK, listNSGsResponse),
				),
			)
			response, err = client.Get("/network_security_groups")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_security_group+json;type=collection"))
		})

		It("lists all network_security_groups]", func() {
			networkSecurityGroups := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNSGsResponse), &networkSecurityGroups)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkSecurityGroups["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-1/network_security_groups")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one network security groups", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi1"),
					ghttp.RespondWith(http.StatusOK, listOneNSGResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/network_security_groups/khrvi1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_security_group+json"))
		})

		It("retrieves an existing availability set", func() {
			var networkSecurityGroups map[string]interface{}
			err := json.Unmarshal([]byte(listOneNSGResponse), &networkSecurityGroups)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkSecurityGroups)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-1/network_security_groups/khrvi")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: khrvi\"}\n"))
		})
	})

	Describe("creating", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi1"),
					ghttp.RespondWith(201, listOneNSGResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-1/network_security_groups", "{\"name\": \"khrvi1\", \"location\": \"westus\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource network security group href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-1/network_security_groups/khrvi1"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi1"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-1/network_security_groups/khrvi1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 204 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(204))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})
})
