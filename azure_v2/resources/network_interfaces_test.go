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
	listNIsResponse   = `{"value":[{"etag":"W/\"854ed6c8-3337-4617-a5a8-38aa5051957b\"","href":"/resource_groups/Group-1/network_interfaces/1_khrvi","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/1_khrvi","location":"westus","name":"1_khrvi","properties":{"dnsSettings":{"dnsServers":["10.1.0.5"]},"ipConfigurations":[{"etag":"W/\"854ed6c8-3337-4617-a5a8-38aa5051957b\"","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/1_khrvi/ipConfigurations/1_khrvi_ip","name":"1_khrvi_ip","properties":{"privateIPAddress":"10.0.0.106","privateIPAllocationMethod":"Dynamic","provisioningState":"Succeeded","subnet":{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/virtualNetworks/khrvi/subnets/khrvi"}}}],"provisioningState":"Succeeded"}}]}`
	listOneNIResponse = `{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/1_khrvi","name":"1_khrvi","location":"westus","etag":"W/\"854ed6c8-3337-4617-a5a8-38aa5051957b\"","properties":{"dnsSettings":{"dnsServers":["10.1.0.5"]},"ipConfigurations":[{"etag":"W/\"854ed6c8-3337-4617-a5a8-38aa5051957b\"","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/1_khrvi/ipConfigurations/1_khrvi_ip","name":"1_khrvi_ip","properties":{"privateIPAddress":"10.0.0.106","privateIPAllocationMethod":"Dynamic","provisioningState":"Succeeded","subnet":{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/virtualNetworks/khrvi/subnets/khrvi"}}}],"provisioningState":"Succeeded"},"href":"/resource_groups/Group-1/network_interfaces/1_khrvi"}`
)

var _ = Describe("network interfaces", func() {

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
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkInterfacePath),
					ghttp.RespondWith(http.StatusOK, listNIsResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/network_interfaces")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_interface+json;type=collection"))
		})

		It("lists all network_interfaces inside one resource group", func() {
			networkInterfaces := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNIsResponse), &networkInterfaces)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkInterfaces["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			subscriptionID := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/"+networkInterfacePath),
					ghttp.RespondWith(http.StatusOK, listNIsResponse),
				),
			)
			response, err = client.Get("/network_interfaces")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_interface+json;type=collection"))
		})

		It("lists all network_interfaces", func() {
			networkInterfaces := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNIsResponse), &networkInterfaces)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkInterfaces["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkInterfacePath),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-1/network_interfaces")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one network interface", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkInterfacePath+"/1_khrvi"),
					ghttp.RespondWith(http.StatusOK, listOneNIResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/network_interfaces/1_khrvi")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network_interface+json"))
		})

		It("retrieves an existing network interface", func() {
			var networkInterface map[string]interface{}
			err := json.Unmarshal([]byte(listOneNIResponse), &networkInterface)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networkInterface)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkInterfacePath+"/khrvi1"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-1/network_interfaces/khrvi1")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: khrvi1\"}\n"))
		})
	})

	Describe("creating", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkInterfacePath+"/1_khrvi"),
					ghttp.RespondWith(201, listOneNIResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-1/network_interfaces", "{\"name\": \"1_khrvi\", \"location\": \"westus\", \"subnet_id\": \"subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/virtualNetworks/khrvi/subnets/khrvi\", \"network_security_group_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkSecurityGroups/khrvi1\", \"public_ip_address_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/publicIPAddresses/khrvi-3\", \"private_ip_address\": \"10.0.0.130\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource network interface href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-1/network_interfaces/1_khrvi"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+networkInterfacePath+"/1_khrvi"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-1/network_interfaces/1_khrvi")
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
