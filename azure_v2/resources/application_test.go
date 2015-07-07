package resources

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
	"github.com/rightscale/self-service-plugins/azure_v2/config"
	am "github.com/rightscale/self-service-plugins/azure_v2/middleware"
)

const (
	authRsponse              = `{"access_token": "test_access_token", "expires_on": "123456789"}`
	servicePrincipalResponse = `{"value":[{ "objectType":"ServicePrincipal","objectId":"7f06b355-4136-4d8b-a3c8-7028f59869ae"}]}`
	listRoleAssignments      = `{"value":[{"properties":{"roleDefinitionId":"/subscriptions/test/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c","principalId":"7f06b355-4136-4d8b-a3c8-7028f59869ae","scope":"/subscriptions/test"},"id":"/subscriptions/test/providers/Microsoft.Authorization/roleAssignments/4f87261d-2816-465d-8311-70a27558df4c","type":"Microsoft.Authorization/roleAssignments","name":"4f87261d-2816-465d-8311-70a27558df4c"}]}`
	deleteRoleAssignment     = "{\"properties\":{\"roleDefinitionId\":\"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c\",\"principalId\":\"7f06b355-4136-4d8b-a3c8-7028f59869ae\",\"scope\":\"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a\"},\"id\":\"/subscriptions/2d2b2267-ff0a-46d3-9912-8577acb18a0a/providers/Microsoft.Authorization/roleAssignments/4f87261d-2816-465d-8311-70a27558df4c\",\"type\":\"Microsoft.Authorization/roleAssignments\",\"name\":\"4f87261d-2816-465d-8311-70a27558df4c\"}"
)

var _ = Describe("application", func() {

	var do *ghttp.Server
	var client *AzureClient
	var response *Response
	var err error

	BeforeEach(func() {
		do = ghttp.NewServer()
		config.AuthHost = do.URL()
		config.GraphURL = do.URL()
		config.BaseURL = do.URL()
		client = NewAzureClient()
	})

	AfterEach(func() {
		do.Close()
	})

	Describe("register", func() {
		BeforeEach(func() {
			// make empty access token and subscription in order to refresh access token
			AccessTokenTest = ""
			*config.SubscriptionIDCred = ""
			CredsTest = am.Credentials{
				TenantID:     "test_tenant",
				ClientID:     "test_client",
				ClientSecret: "test_secret",
				RefreshToken: "test_token",
				Subscription: "test_subscription",
			}
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/test_tenant/oauth2/token"),
					ghttp.RespondWith(http.StatusOK, authRsponse),
				),
				// requesting access token with resource "https://graph.windows.net/"
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/test_tenant/oauth2/token"),
					ghttp.RespondWith(http.StatusOK, authRsponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/test_tenant/servicePrincipals", "api-version=1.5&$filter=appId%20eq%20'test_client'"),
					ghttp.RespondWith(http.StatusOK, servicePrincipalResponse),
				),

				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", MatchRegexp(`/subscriptions/test_subscription/providers/microsoft.authorization/roleassignments/\w`)),
					ghttp.RespondWith(201, ""),
				),
			)
			//TODO: try to use analog of before_all
			response, err = client.Post("/application/register", "")
		})

		AfterEach(func() {
			AccessTokenTest = "fake"                    // set back access token
			*config.SubscriptionIDCred = subscriptionID //set back default subscription
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 201 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(4))
			Ω(response.Status).Should(Equal(201))
		})

		It("returns empty response in case of 201", func() {
			Ω(response.Body).Should(BeEmpty())
		})
		It("returns Access Token in the cookie", func() {
			Ω(response.Cookies[0].Name).Should(Equal("AccessToken"))
			Ω(response.Cookies[0].Value).Should(Equal("test_access_token"))
		})
		It("returns ExpiresOn in the cookie", func() {
			Ω(response.Cookies[1].Name).Should(Equal("ExpiresOn"))
			Ω(response.Cookies[1].Value).Should(Equal("123456789"))
		})
	})

	Describe("unregister", func() {
		BeforeEach(func() {
			do.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/test_tenant/oauth2/token"),
					ghttp.RespondWith(http.StatusOK, authRsponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/test_tenant/servicePrincipals", "api-version=1.5&$filter=appId%20eq%20'test_client'"),
					ghttp.RespondWith(http.StatusOK, servicePrincipalResponse),
				),

				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/subscriptions/test/providers/microsoft.authorization/roleassignments"),
					ghttp.RespondWith(200, listRoleAssignments),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/subscriptions/test/providers/microsoft.authorization/roleassignments/4f87261d-2816-465d-8311-70a27558df4c"),
					ghttp.RespondWith(200, deleteRoleAssignment),
				),
			)
			//TODO: try to use analog of before_all
			response, err = client.Delete("/application/unregister")
		})

		It("no error occured", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns 204 status code", func() {
			Ω(do.ReceivedRequests()).Should(HaveLen(4))
			Ω(response.Status).Should(Equal(204))
		})

		It("returns empty response", func() {
			Ω(response.Body).Should(BeEmpty())
		})
	})
})
