namespace "dme" do
  service do
    host "dme.test.rightscale.com"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/dme/accounts/:account_id"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service" ,     # special headers as needed
      "X-Api-Version" => "1.0",
      "X-Api-Shared-Secret" => "7uNzca10X8&y$11L&0OM9fLDkLIma*q&P9jJeZG@#Gf"
    } end
  end
  type "record" do
    provision "provision_record"
    delete "delete_record"
    fields do
      field "domain" do
        type "string"
        required true
      end
      field "name" do
        type "string"
        required true
      end
      field "value" do
        type "string"
        required true
      end
      field "type" do
        type "string"
        required true
      end
      field "dynamicDns" do
        type "boolean"
      end
      field "ttl" do
        type "number"
      end
    end
  end
end 