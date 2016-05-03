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
	listInstancesEmptyResponse = `{"value":[]}`
	listInstancesResponse      = `{"value":[{"href":"/resource_groups/Group-1/instances/khrvi","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Compute/virtualMachines/khrvi","location":"westus","name":"khrvi","properties":{"hardwareProfile":{"vmSize":"Standard_G1"},"networkProfile":{"networkInterfaces":[{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni"}]},"provisioningState":"failed","storageProfile":{"dataDisks":[],"osDisk":{"caching":"ReadWrite","name":"os-asdasdasda-rs","osType":"Linux","vhd":{"uri":"https://khrvitestgo.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"}}}},"type":"Microsoft.Compute/virtualMachines"}]}`
	listOneInstanceResponse    = `{"href":"/resource_groups/Group-1/instances/khrvi","id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Compute/virtualMachines/khrvi","location":"westus","name":"khrvi","properties":{"hardwareProfile":{"vmSize":"Standard_G1"},"networkProfile":{"networkInterfaces":[{"id":"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni"}]},"provisioningState":"failed","storageProfile":{"dataDisks":[],"osDisk":{"caching":"ReadWrite","name":"os-asdasdasda-rs","osType":"Linux","vhd":{"uri":"https://khrvitestgo.blob.core.windows.net/vhds/khrvi_image-os-2015-05-18.vhd"}}}},"type":"Microsoft.Compute/virtualMachines"}`
	recordNotFound             = `{"error":{"code":"ResourceNotFound","message":"Resource not found."}}`
)

var _ = Describe("instances", func() {

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
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/instances")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.instance+json;type=collection"))
		})

		It("lists all instances inside one resource group", func() {
			instances := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listInstancesResponse), &instances)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(instances["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing via 'flat' route", func() {
		BeforeEach(func() {
			subscriptionID := "test"
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/"+virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesResponse),
				),
			)
			response, err = client.Get("/instances")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.instance+json;type=collection"))
		})

		It("lists all instances inside one resource group", func() {
			instances := make(map[string]interface{}, 0)
			err = json.Unmarshal([]byte(listInstancesResponse), &instances)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(instances["value"])
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("listing empty", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath),
					ghttp.RespondWith(http.StatusOK, listInstancesEmptyResponse),
				),
			)
		})

		It("returns empty array", func() {
			response, err = client.Get("/resource_groups/Group-1/instances")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
			Ω(response.Body).Should(Equal("[]\n"))
		})
	})

	Describe("list one instance", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi"),
					ghttp.RespondWith(http.StatusOK, listOneInstanceResponse),
				),
			)
			response, err = client.Get("/resource_groups/Group-1/instances/khrvi")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 200 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(200))
		})

		It("returns a resource specific content type in the header", func() {
			Ω(response.Headers["Content-Type"][0]).Should(Equal("vnd.rightscale.instance+json"))
		})

		It("retrieves an existing instance", func() {
			var instance map[string]interface{}
			err := json.Unmarshal([]byte(listOneInstanceResponse), &instance)
			Expect(err).NotTo(HaveOccurred())
			expected, err := json.Marshal(instance)
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Body).Should(MatchJSON(expected))
		})
	})

	Describe("retrieving a non-existant resource", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi1"),
					ghttp.RespondWith(http.StatusNotFound, recordNotFound),
				),
			)
		})

		It("returns 404", func() {
			response, err = client.Get("/resource_groups/Group-1/instances/khrvi1")
			Expect(err).NotTo(HaveOccurred())
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(404))
			Ω(response.Body).Should(Equal("{\"Code\":404,\"Message\":\"Could not find resource with id: khrvi1\"}\n"))
		})
	})

	Describe("creating with public image", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi"),
					ghttp.VerifyJSONRepresenting(requestParams{
						Name:     "khrvi",
						Location: "westus",
						Properties: map[string]interface{}{
							"hardwareProfile": map[string]interface{}{"vmSize": "Standard_G1"},
							"storageProfile": map[string]interface{}{
								"imageReference": map[string]interface{}{
									"publisher": "a10networks",
									"offer":     "a10-vthunder-adc",
									"sku":       "vthunder_100mbps",
									"version":   "1.0.0",
								},
								"osDisk": map[string]interface{}{
									"name":         "os-khrvi-rs",
									"caching":      "ReadWrite",
									"createOption": "FromImage",
									"vhd": map[string]interface{}{
										"uri": "https://khrvitestgo1.blob.core.windows.net/vhds/os-khrvi-rs.vhd",
									},
								},
							},
							"osProfile": map[string]interface{}{
								"computerName":  "khrvi",
								"adminUsername": "rsadministrator",
								"adminPassword": "Pass1234@",
							},
							"networkProfile": map[string]interface{}{
								"networkInterfaces": []map[string]interface{}{
									{"id": "/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni"},
								},
							},
						}}),
					ghttp.RespondWith(201, listOneInstanceResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"instance_type_uid\": \"Standard_G1\", \"location\": \"westus\", \"network_interface_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni\", \"image_id\": \"/Subscriptions/test/Providers/Microsoft.Compute/Locations/westus/Publishers/a10networks/ArtifactTypes/VMImage/Offers/a10-vthunder-adc/Skus/vthunder_100mbps/Versions/1.0.0\", \"storage_account_id\": \"/subscriptions/test/resourceGroups/group-1/providers/Microsoft.Storage/storageAccounts/khrvitestgo1\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource instance href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-1/instances/khrvi"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("creating with user image and user_data", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi"),
					ghttp.VerifyJSONRepresenting(requestParams{
						Name:     "khrvi",
						Location: "westus",
						Properties: map[string]interface{}{
							"hardwareProfile": map[string]interface{}{"vmSize": "Standard_G1"},
							"storageProfile": map[string]interface{}{
								"osDisk": map[string]interface{}{
									"name":         "os-khrvi-rs",
									"caching":      "ReadWrite",
									"createOption": "FromImage",
									"vhd": map[string]interface{}{
										"uri": "https://khrvitestgo1.blob.core.windows.net/vhds/os-khrvi-rs.vhd",
									},
									"osType": "Linux",
									"image": map[string]interface{}{
										"uri": "https://khrvitesttest1.blob.core.windows.net/vhds/os-khrvi-rs.vhd",
									},
								},
							},
							"osProfile": map[string]interface{}{
								"computerName":  "khrvi",
								"adminUsername": "rsadministrator",
								"adminPassword": "Pass1234@",
								"customData":    "dGVzdF91c2VyX2RhdGE=",
							},
							"networkProfile": map[string]interface{}{
								"networkInterfaces": []map[string]interface{}{
									{"id": "/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni"},
								},
							},
						}}),
					ghttp.RespondWith(201, listOneInstanceResponse),
				),
			)
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"user_data\":\"test_user_data\", \"instance_type_uid\": \"Standard_G1\", \"location\": \"westus\", \"network_interface_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni\", \"private_image_id\": \"https://khrvitesttest1.blob.core.windows.net/vhds/os-khrvi-rs.vhd\", \"storage_account_id\": \"/subscriptions/test/resourceGroups/group-1/providers/Microsoft.Storage/storageAccounts/khrvitestgo1\"}")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(1))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns a resource instance href in the 'Location' header", func() {
			Ω(response.Headers["Location"][0]).Should(Equal("/resource_groups/Group-1/instances/khrvi"))
		})

		It("return empty body", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})

	Describe("creating with wrong params", func() {
		It("returns validation error about missing 'name'", func() {
			response, err = client.Post("/resource_groups/Group-1/instances", "{}")
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Status).Should(Equal(400))
			Ω(response.Body).Should(Equal("{\"Code\":400,\"Message\":\"You have specified an invalid 'name' parameter.\"}\n"))
		})

		It("returns validation error about missing 'location'", func() {
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"instance_type_uid\": \"Standard_G1\", \"network_interface_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni\", \"image_id\": \"/Subscriptions/test/Providers/Microsoft.Compute/Locations/westus/Publishers/a10networks/ArtifactTypes/VMImage/Offers/a10-vthunder-adc/Skus/vthunder_100mbps/Versions/1.0.0\", \"storage_account_id\": \"/subscriptions/test/resourceGroups/group-1/providers/Microsoft.Storage/storageAccounts/khrvitestgo1\"}")
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Status).Should(Equal(400))
			Ω(response.Body).Should(Equal("{\"Code\":400,\"Message\":\"You have specified an invalid 'location' parameter.\"}\n"))
		})

		It("returns validation error about missing 'image_id'", func() {
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"location\": \"westus\", \"instance_type_uid\": \"Standard_G1\", \"network_interface_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni\", \"storage_account_id\": \"/subscriptions/test/resourceGroups/group-1/providers/Microsoft.Storage/storageAccounts/khrvitestgo1\"}")
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Status).Should(Equal(400))
			Ω(response.Body).Should(Equal("{\"Code\":400,\"Message\":\"One of these two params should be passed: 'image_id' or 'private_image_id'.\"}\n"))
		})

		It("returns validation error about wrong 'image_id'", func() {
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"location\": \"westus\", \"image_id\": \"/Subscriptions/test/Providers/Microsoft.Compute/Locations/westus/Publishers/a10networks/ArtifactTypes/VMImage/Offers/a10-vthunder-adc/Skus/vthunder_100mbps\", \"instance_type_uid\": \"Standard_G1\", \"network_interface_id\": \"/subscriptions/test/resourceGroups/Group-1/providers/Microsoft.Network/networkInterfaces/khrvi_ni\", \"storage_account_id\": \"/subscriptions/test/resourceGroups/group-1/providers/Microsoft.Storage/storageAccounts/khrvitestgo1\"}")
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Status).Should(Equal(400))
			Ω(response.Body).Should(Equal("{\"Code\":400,\"Message\":\"You have specified an invalid 'image_id' parameter.\"}\n"))
		})

		It("returns validation error about missing 'storage_account_id'", func() {
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"location\": \"westus\", \"instance_type_uid\": \"Standard_G1\", \"image_id\": \"/Subscriptions/test/Providers/Microsoft.Compute/Locations/westus/Publishers/a10networks/ArtifactTypes/VMImage/Offers/a10-vthunder-adc/Skus/vthunder_100mbps/Versions/1.0.0\"}")
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Status).Should(Equal(400))
			Ω(response.Body).Should(Equal("{\"Code\":400,\"Message\":\"You have specified an invalid 'storage_account_id' parameter.\"}\n"))
		})

		It("returns validation error about missing 'instance_type_id'", func() {
			response, err = client.Post("/resource_groups/Group-1/instances", "{\"name\": \"khrvi\", \"location\": \"westus\", \"image_id\": \"/Subscriptions/test/Providers/Microsoft.Compute/Locations/westus/Publishers/a10networks/ArtifactTypes/VMImage/Offers/a10-vthunder-adc/Skus/vthunder_100mbps/Versions/1.0.0\", \"storage_account_id\": \"/subscriptions/test/resourceGroups/group-1/providers/Microsoft.Storage/storageAccounts/khrvitestgo1\"}")
			Expect(err).NotTo(HaveOccurred())
			Ω(response.Status).Should(Equal(400))
			Ω(response.Body).Should(Equal("{\"Code\":400,\"Message\":\"You have specified an invalid 'instance_type_id' parameter.\"}\n"))
		})
	})

	Describe("deleting", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/"+subscriptionID+"/resourceGroups/Group-1/"+virtualMachinesPath+"/khrvi"),
					ghttp.RespondWith(200, ""),
				),
			)
			response, err = client.Delete("/resource_groups/Group-1/instances/khrvi")
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
