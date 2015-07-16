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
	listNetworksResponse   = `{"value":[{"etag":"W/\"b055718a-6d32-49e4-a4dd-d8bde3f84070\"","href":"/resource_groups/Group-3/networks/khrvi-3","id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Network/virtualNetworks/khrvi-3","location":"westus","name":"khrvi-3","properties":{"addressSpace":{"addressPrefixes":["10.0.0.0/16"]},"provisioningState":"Succeeded","subnets":[{"etag":"W/\"b055718a-6d32-49e4-a4dd-d8bde3f84070\"","id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Network/virtualNetworks/khrvi-3/subnets/khrvi-3","name":"khrvi-3","properties":{"addressPrefix":"10.0.0.0/16","provisioningState":"Succeeded"}}]}},{"etag":"W/\"2bc1c8a9-e9d1-4432-8b92-8e6c79d48e82\"","href":"/resource_groups/Group-3/networks/net2","id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Network/virtualNetworks/net2","location":"westus","name":"net2","properties":{"addressSpace":{"addressPrefixes":["10.0.0.0/16"]},"dhcpOptions":{"dnsServers":["10.1.0.5","10.1.0.6"]},"provisioningState":"Succeeded"}}]}`
	listOneNetworkResponse = `{"id":"/subscriptions/test/resourceGroups/Group-3/providers/Microsoft.Network/virtualNetworks/net2","name":"net2","location":"westus","etag":"W/\"2bc1c8a9-e9d1-4432-8b92-8e6c79d48e82\"","properties":{"addressSpace":{"addressPrefixes":["10.0.0.0/16"]},"dhcpOptions":{"dnsServers":["10.1.0.5","10.1.0.6"]},"provisioningState":"Succeeded"},"href":"/resource_groups/Group-3/networks/net2"}`
)

var _ = Describe("networks", func() {

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
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath),
					ghttp.RespondWith(http.StatusOK, listNetworksResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-3/networks")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network+json;type=collection"))
		})

		It("lists all networks inside one resource group", func() {
			networks := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNetworksResponse), &networks)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networks["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/"+networkPath),
					ghttp.RespondWith(http.StatusOK, listNetworksResponse),
				),
			)
			response, err = client.Get("/networks")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network+json;type=collection"))
		})

		It("lists all networks inside one resource group", func() {
			networks := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listNetworksResponse), &networks)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(networks["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath),
					ghttp.RespondWith(http.StatusOK, listEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-3/networks")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one network", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/net1"),
					ghttp.RespondWith(http.StatusOK, listOneNetworkResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-3/networks/net1")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.network+json"))
		})

		It("retrieves an existing network", func() {
			var network map[string]interface{}
			err := json.Unmarshal([]byte(listOneNetworkResponse), &network)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(network)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/khrvi"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-3/networks/khrvi")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: khrvi\"}\n"))
		})
	})

	Describe("creating with no subnets", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/net2"),
					ghttp.VerifyJSONRepresenting(networkRequestParams{
						Name:     "net2",
						Location: "westus",
						Properties: map[string]interface{}{
							"addressSpace": map[string]interface{}{
								"addressPrefixes": []string{"10.0.0.0/16"},
							},
							"subnets": nil,
						},
					}),
					ghttp.RespondWith(201, listOneNetworkResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-3/networks", "{\"name\": \"net2\", \"location\": \"westus\", \"address_prefixes\": [\"10.0.0.0/16\"]}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource network href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-3/networks/net2"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("creating with one subnet", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/net2"),
					ghttp.VerifyJSONRepresenting(networkRequestParams{
						Name:     "net2",
						Location: "westus",
						Properties: map[string]interface{}{
							"addressSpace": map[string]interface{}{
								"addressPrefixes": []string{"10.0.0.0/16"},
							},
							"subnets": []map[string]interface{}{
								{"name": "subnet_name",
									"properties": map[string]interface{}{
										"addressPrefix": "10.0.0.0/16",
									},
								},
							},
							"dhcpOptions": map[string]interface{}{
								"dnsServers": []string{"10.1.0.5", "10.1.0.6"},
							},
						},
					}),
					ghttp.RespondWith(201, listOneNetworkResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-3/networks", "{\"name\": \"net2\", \"location\": \"westus\", \"address_prefixes\": [\"10.0.0.0/16\"], \"subnets\": [{\"name\": \"subnet_name\", \"address_prefix\": \"10.0.0.0/16\"}], \"dhcp_options\": {\"dnsServers\": [\"10.1.0.5\", \"10.1.0.6\"]}}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource network href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-3/networks/net2"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-3/"+networkPath+"/net1"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-3/networks/net1")
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
