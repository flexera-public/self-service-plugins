name "Configurator - Single Chef Node"
rs_ca_ver 20131202
short_description "Configurator - Single Chef Node"
long_description "This CAT uses the configurator plugin to launch a raw image as a RightScale server configured by a pre-existing Chef installation"

###########
# Namespace
###########

# For clarity sake
::IP_REGEXP = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
::HOST_REGEXP = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$"

namespace "cm" do
  service do
    host  "54.184.12.120" #"cm.test.rightscale.com"
    path "/cm/accounts/:account_id"
    headers do {
      "X-Api-Version" => "1.0",
      "X-Secret" => "R9Bt4GMqQBT3UZREoW7aaACUnWLWdGqO"
    } end
  end

  type "configuration" do
    provision "provision_configuration"
    delete "delete_configuration"
    fields do
      field "type" do
        type "string"
        regexp "^chef$" # Only chef for now
        required true
      end
      field "settings" do
        type "composite"
        required true
      end
    end
  end

  type "booter" do
    provision "provision_booter"
    delete "delete_booter"
    fields do
      field "host" do
        type "string"
        required "true"
        regexp "(?:#{::IP_REGEXP}|#{::HOST_REGEXP})"
      end
      field "ssh_key" do
        type "resource"
      end
    end
  end
end

# Define the RCL definitions to create and destroy the resource
define provision_configuration(@raw_conf) return @conf do
  $obj = to_object(@raw_conf)
  $fields = $obj["fields"]
  @conf = cm.configuration.create($fields) # Calls .create on the API resource
end

define delete_configuration(@conf) do
  @conf.destroy() # Calls .delete on the API resource
end

define provision_booter(@raw_booter) return @booter do
  $obj = to_object(@raw_booter)
  $fields = $obj["fields"]
  @conf = cm.booter.create($fields) # Calls .create on the API resource
end

define delete_booter(@booter) do
  @booter.destroy() # Calls .delete on the API resource
end

#########
# Parameters
#########
parameter "chef_server_url" do
  type "string"
  label "Chef Server URL"
  category "Chef"
  default "https://api.opscode.com/organizations/rs-st-dev"
end

parameter "validation_key" do
  type "string"
  label "Chef Validation Key"
  description "Name of RightScale credential holding the Chef server validation key"
  default "rs-st-dev-validator"
end

parameter "environment" do
  type "string"
  label "Chef environment"
  category "Chef"
  default "_default"
end

parameter "run_list" do
  type "list"
  label "Boot run list"
  category "Chef"
end

parameter "first_attributes" do
  type "string"
  label "Attributes used to run initial configuration"
  category "Chef"
end

#########
# Resources
#########
resource "chef_cm", type: "cm.configuration" do
  type "chef"
  settings do {
    "chef_server_url" => $chef_server_url,
    "validation_key_name" => $validation_key,
    "run_list" => $run_list,
    "chef_environment" => $environment,
    "first_attributes" => $first_attributes
  } end
end

resource "cm_server", type: "server" do
  name "cm_server"
  cloud_href "/api/clouds/1"
  instance_type "m1.small"
  ssh_key "default"
  user_data "@chef_cm.bootstrap_script" # server must be tagged with 'rs_agent:userdata=mime'
  server_template find('RL10.0.rc2 Linux Base') # Could be anything
end
