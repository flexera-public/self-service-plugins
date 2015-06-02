name "Base Linux for Azure"
rs_ca_ver 20131202
short_description "Allows you to launch a machine on Azure"

namespace "azure" do
  service do
    host "cmdev-selfservice-403.test.rightscale.com"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/azure_plugin"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service"      # special headers as needed
    } end
  end
  type "instance" do                           # defines resource of type "pod"
    provision "provision_instance"             # name of RCL definition to use to provision the resource
    delete "delete_instance"                   # name of RCL definition to use to delete the resource
    fields do
      field "name" do                               
        type "string"
        required true
      end
      field "location" do                               
        type "string"
        required true
      end
      field "instance_type_uid" do                               
        type "string"
        required true
      end
      field "group_name" do                               
        type "string"
        required true
      end
    end
  end
end

parameter "instance_name" do
  type "string"
  label "Instance Name"
  description "The name to give the instance"
end

parameter "instance_size" do
  type "string"
  label "Instance Size"
  description "Size of the instance"
  default "Standard_G1"
  allowed_values "Standard_G1", "Standard_G2", "Standard_G3"
end

resource "base_server", type: "azure.instance" do
  name                  $instance_name
  instance_type_uid     $instance_size
  location              "westus"
  group_name            "Group-1"
end
 
# Define the RCL definitions to create and destroy the resource
define provision_instance(@raw_instance) return @instance do
  $obj = to_object(@raw_instance)
  $fields = $obj["fields"]
  @instance = azure.instance.create($fields) # Calls .create on the API resource
end
define delete_instance(@instance) do
  @instance.destroy() # Calls .delete on the API resource
end