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
	listNSGRsResponse   = `{"value":[{"etag":"W/\"f16517fd-a1f6-449d-a77c-06ef63318423\"","href":"/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules/khrvi2_1","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/securityRules/khrvi2_1","name":"khrvi2_1","properties":{"access":"Allow","description":"test","destinationAddressPrefix":"*","destinationPortRange":"800","direction":"Inbound","priority":100,"protocol":"Tcp","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"800"}}]}`
	listOneNSGRResponse = `{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi2/securityRules/khrvi2_1","name":"khrvi2_1","etag":"W/\"f16517fd-a1f6-449d-a77c-06ef63318423\"","properties":{"access":"Allow","description":"test","destinationAddressPrefix":"*","destinationPortRange":"800","direction":"Inbound","priority":100,"protocol":"Tcp","provisioningState":"Succeeded","sourceAddressPrefix":"*","sourcePortRange":"800"},"href":"/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules/khrvi2_1"}`
)

var _ = Describe("network security group rules", func() {

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
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi2/securityRules"),
					ghttp.RespondWith(http.StatusOK, listNSGRsResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_security_group_rule+json;type=collection"))
		})

		It("lists all network security group rule inside one resource group", func() {
			networkSecurityGroupRules := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNSGRsResponse), &networkSecurityGroupRules)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkSecurityGroupRules["value"])
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
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi1/securityRules"),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi2/securityRules"),
					ghttp.RespondWith(http.StatusOK, listNSGRsResponse),
				),
			)
			response, err = client.Get("/network_security_group_rules")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(3))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_security_group_rule+json;type=collection"))
		})

		It("lists all network security group rule", func() {
			networkSecurityGroupRules := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNSGRsResponse), &networkSecurityGroupRules)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkSecurityGroupRules["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkSecurityGroupPath+"/khrvi2/securityRules"),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-3/network_security_groups/khrvi2/network_security_group_rules")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one network security group rule", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi2/securityRules/khrvi2_1"),
					ghttp.RespondWith(http.StatusOK, listOneNSGRResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules/khrvi2_1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_security_group_rule+json"))
		})

		It("retrieves an existing network security group rule", func() {
			var networkSecurityGroupRule map[string]interface{}
			err := json.Unmarshal([]byte(listOneNSGRResponse), &networkSecurityGroupRule)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkSecurityGroupRule)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi2/securityRules/khrvi"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules/khrvi")
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
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi2/securityRules/khrvi2_1"),
					ghttp.RespondWith(201, listOneNSGRResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules", "{\"name\": \"khrvi2_1\", \"description\": \"test\", \"protocol\": \"Tcp\", \"source_port_range\": \"801\", \"destination_port_range\": \"801\", \"source_address_prefix\": \"*\", \"destination_address_prefix\": \"*\", \"access\": \"Allow\", \"priority\": 200, \"direction\": \"Inbound\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource network security group rule href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules/khrvi2_1"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkSecurityGroupPath+"/khrvi2/securityRules/khrvi2_1"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-1/network_security_groups/khrvi2/network_security_group_rules/khrvi2_1")
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
