name "Test for CloudFormation Wordpress Template"
rs_ca_ver 20131202
short_description "Test for a CloudFormation Template

![logo](http://d2wwfe3odivqm9.cloudfront.net/wp-content/uploads/2014/05/aws_icon-cloudformation_white-200x200.png)"


output 'dns_name' do
  label "URL"
  category "General"
  default_value join(["http://",$domain_name,".dev.rightscaleit.com/wordpress"])
end

# Cloud Selection
parameter "domain_name" do
  type "string"
  label "Name for A record (will be appended with .dev.rightscaleit.com)"
  category "DNS"
end

resource 'dns_entry', type: 'dme.record' do
  domain 'dev.rightscaleit.com' 
  name $domain_name
  type "A"
  value "192.192.192.192"
  ttl 30
end

parameter "cft_param_instance_type" do
  type "string"
  label "Instance Type"
  default "m1.small"
  allowed_values "t2.medium", "m1.small", "m1.medium", "m1.large", "m1.xlarge", "m2.xlarge", "m2.2xlarge"
end

parameter "cft_param_key" do
  type "string"
  label "Key to use for the instance"
  default "default"
end

parameter "cft_dbrootpw" do
  label "DBRootPassword"
  type "string"
  allowed_pattern "[a-zA-Z0-9]*"
  min_length 8
  max_length 41
  description "MySQL root password"
  constraint_description "must contain only alphanumeric characters."
  default "asdlk9u43iou"
end

parameter "cft_dbuser" do
  label "DBUser"
  type "string"
  allowed_pattern "[a-zA-Z][a-zA-Z0-9]*"
  min_length 1
  max_length 16
  description "The WordPress database admin account username"
  constraint_description "must begin with a letter and contain only alphanumeric characters."
  default "sqluser"
end


parameter "cft_dbpw" do
  label "DBPassword"
  type "string"
  allowed_pattern "[a-zA-Z0-9]*"
  min_length 8
  max_length 41
  description "The WordPress database admin account password"
  constraint_description "must contain only alphanumeric characters."
  default "sdfaoi9u0fdsa"
end

resource 'cft_template', type: 'ec2cft.stack' do
  name @@deployment.name
  template "https://s3.amazonaws.com/rol-test-cf-templates/WordPress_Single_Instance_CAT.template"
  parameters do {
    "KeyName" => $cft_param_key,
    "InstanceType" => $cft_param_instance_type,
    "DBRootPassword" => $cft_dbrootpw,
    "DBUser" => $cft_dbuser,
    "DBPassword" => $cft_dbpw
  } end
end

operation 'enable' do
  definition 'update_ip'
  description 'Update the IP of the record'
end

define update_ip(@dns_entry, @cft_template) return @dns_entry do
  $ips = select(@cft_template.outputs, { "key": "InstanceIP" })
  if size($ips) > 0     
    @dns_entry.update(name: @dns_entry.name, type: @dns_entry.type, value: $ips[0]["value"])
  end
end

############################
############################
 #   ___ _             _ ___                   _   _          
 #  / __| |___ _  _ __| | __|__ _ _ _ __  __ _| |_(_)___ _ _  
 # | (__| / _ \ || / _` | _/ _ \ '_| '  \/ _` |  _| / _ \ ' \ 
 #  \___|_\___/\_,_\__,_|_|\___/_| |_|_|_\__,_|\__|_\___/_||_|
 #                                                           
 #
############################
############################
               
namespace "ec2cft" do
  service do
    host "http://52.12.86.212:8082"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/ec2cft/accounts/:account_id"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service" ,     # special headers as needed
      "X-API-Version" => "1.0",
      "X-Api-Shared-Secret" => ""      
    } end
  end
  type "stack" do                          
    provision "provision_stack"            
    delete "delete_stack"                  
    fields do
      field "template" do
        type "string"
        required true
      end
      field "name" do
        type "string"
        required true
      end
      field "parameters" do
        type "string"
      end
    end
  end
end 

define provision_stack(@raw_stack) return @resource do
  $obj = to_object(@raw_stack)
  $to_create = $obj["fields"]
  @resource = ec2cft.stack.create($to_create)
  sleep_until(@resource.status != "CREATE_IN_PROGRESS")
  if @resource.status != "CREATE_COMPLETE"
    $status = @resource.status
    @resource.destroy()
    raise "Error creating CFT stack: " + $status
  end
end

define delete_stack(@stack) do
  @stack.destroy()
end

############################
############################
 #  ___  __  __ ___
 # |   \|  \/  | __|
 # | |) | |\/| | _|
 # |___/|_|  |_|___|
 #
############################
############################

namespace "dme" do
  service do
    host "http://50.112.115.235:8081" # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/dme/accounts/:account_id"  # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service" ,     # special headers as needed
      "X-Api-Version" => "1.0",
      "X-Api-Shared-Secret" => ""
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

define provision_record(@raw_record) return @resource do
  $obj = to_object(@raw_record)
  $to_create = $obj["fields"]
  @resource = dme.record.create($to_create)
end

define delete_record(@record) do
  @record.destroy()
end

