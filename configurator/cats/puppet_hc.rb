name "Puppet - Single Node"
rs_ca_ver 20131202
short_description "Puppet - Single Chef Node"
long_description "This CAT uses cloud init to bootstrap a puppet agent"

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
resource "chef_cm", type: "cm.chef_configuration" do
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

