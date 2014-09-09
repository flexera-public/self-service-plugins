namespace "dme" do
  service do
    host "dme.test.rightscale.com"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/dme/accounts/:account_id"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service"      # special headers as needed
    } end
  end
  type "record" do                          
    provision "provision_record"            
    delete "delete_record"                  
    # path "/records" # Unneeded since we'll use the name of the type by default
    fields do
      id do                              
        type "number",
        required true
      end
      domain do
        type "string",
        required true
      end
      name do
        type "string",
        required true
      end
      value do
        type "string",
        default "1.1.1.1"
        required true
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