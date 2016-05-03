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
	listIPsResponse   = `{"value":[{"etag":"W/\"55777984-d0cc-4385-8961-32734c66974c\"","href":"/resource_groups/Group-1/ip_addresses/khrvi-3","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/publicIPAddresses/khrvi-3","location":"westus","name":"khrvi-3","properties":{"idleTimeoutInMinutes":4,"provisioningState":"Succeeded","publicIPAllocationMethod":"Dynamic"}},{"etag":"W/\"0ce16543-97dc-4e2c-b26c-2f117e8d0742\"","href":"/resource_groups/Group-1/ip_addresses/khrvi_test_static","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/publicIPAddresses/khrvi_test_static","location":"westus","name":"khrvi_test_static","properties":{"idleTimeoutInMinutes":4,"ipAddress":"138.91.196.80","provisioningState":"Succeeded","publicIPAllocationMethod":"Static"}}]}`
	listOneIPResponse = `{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/publicIPAddresses/khrvi_test_static","name":"khrvi_test_static","location":"westus","etag":"W/\"0ce16543-97dc-4e2c-b26c-2f117e8d0742\"","properties":{"idleTimeoutInMinutes":4,"ipAddress":"138.91.196.80","provisioningState":"Succeeded","publicIPAllocationMethod":"Static"},"href":"/resource_groups/Group-1/ip_addresses/khrvi_test_static"}`
)

var _ = Describe("ip addresses", func() {

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
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+ipAddressPath),
					ghttp.RespondWith(http.StatusOK, listIPsResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/ip_addresses")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.ip_address+json;type=collection"))
		})

		It("lists all ip_addresses inside one resource group", func() {
			ipAddresses := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listIPsResponse), &ipAddresses)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(ipAddresses["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			subscriptionID := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/"+ipAddressPath),
					ghttp.RespondWith(http.StatusOK, listIPsResponse),
				),
			)
			response, err = client.Get("/ip_addresses")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.ip_address+json;type=collection"))
		})

		It("lists all ip_addresses inside one resource group", func() {
			ipAddresses := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listIPsResponse), &ipAddresses)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(ipAddresses["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+ipAddressPath),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-1/ip_addresses")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one ip_address", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+ipAddressPath+"/khrvi_test_static"),
					ghttp.RespondWith(http.StatusOK, listOneIPResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/ip_addresses/khrvi_test_static")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.ip_address+json"))
		})

		It("retrieves an existing ip_address", func() {
			var ipAddress map[string]interface{}
			err := json.Unmarshal([]byte(listOneIPResponse), &ipAddress)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(ipAddress)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+ipAddressPath+"/khrvi1"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-1/ip_addresses/khrvi1")
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
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+ipAddressPath+"/khrvi_test_static"),
					ghttp.VerifyJSONRepresenting(ipAddressRequestParams{
						Location: "westus",
						Properties: map[string]interface{}{
							"publicIPAllocationMethod": "Static",
							"idleTimeoutInMinutes":     10,
						},
					}),
					ghttp.RespondWith(201, listOneIPResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-1/ip_addresses", "{\"name\": \"khrvi_test_static\", \"allocation_method\": \"Static\", \"location\": \"westus\", \"timeout\": 10}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource ip_address href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-1/ip_addresses/khrvi_test_static"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+ipAddressPath+"/khrvi_test_static"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-1/ip_addresses/khrvi_test_static")
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
