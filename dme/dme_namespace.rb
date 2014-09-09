namespace "dme" do
  service do
    host "dme.test.rightscale.com"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/dme/accounts/:account_id"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service"      # special headers as needed
    } end
  end
  type "record" do                           # defines resource of type "pod"
    provision "provision_record"             # name of ?
    delete "delete_record"                   # name of ?
    path "/records"
    fields do
      id do                               # field of a pod with rules for validation
        type "number",
        required
      end
      domain do
        type "string",
        required
      end
      name do
        type "string",
        required
      end
      value do
        type "string",
        default "1.1.1.1"
        required
      end
      type do
        type "string"
      end
      dynamicDns do
        type "boolean"
      end
      ttl do
        type "number"
      end
      password do
        type "string"
      end
    end
  end
end 