name "Test for CFT Wordpress w/ DME"
rs_ca_ver 20131202
short_description "Test for a CloudFormation Template with a DME name

![logo](http://d2wwfe3odivqm9.cloudfront.net/wp-content/uploads/2014/05/aws_icon-cloudformation_white-200x200.png)"


output 'dns_name' do
  label "URL"
  category "General"
  default_value join(["http://",$domain_name,".dev.rightscaleit.com/wordpress"])
end

parameter "cft_param_DBRootPassword" do
  type "string"
  no_echo true
  label "Password to use for the root DB login"
  description "Must be at least 6 characters (no spaces)"
  allowed_pattern "^.{6,50}$"
end
# output 'dns_ip' do
#   label "IP Address"
#   category "General"
#   default_value @dns_entry.value
# end

# Cloud Selection
parameter "domain_name" do
  type "string"
  label "Name for A record"
  category "DNS"
end

resource 'dns_entry', type: 'dme.record' do
  domain 'dev.rightscaleit.com' 
  name $domain_name
  type "A"
  value "192.192.192.192"
  ttl 30
end

resource 'cft_template', type: 'ec2cft.stack' do
  name @@deployment.name
  template "https://s3.amazonaws.com/rol-cf-templates/WordPress_Single_Instance_SS.template"
  parameters join(["{\"DBRootPassword\":\"", $cft_param_DBRootPassword, "\",\"KeyName\":\"my normal key\"}"])
end

operation 'enable' do
  definition 'update_ip'
  description 'Update the IP of the record'
end

operation 'terminate' do
  definition 'terminate'
  description 'Terminate and destroy the resources'
end

define terminate(@dns_entry, @cft_template) do
  call delete_stack(@cft_template)
  call delete_record(@dns_entry)
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
    provision "provision_record"            
    delete "delete_record"                  
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

define provision_record(@raw_stack) return @resource do
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

define delete_record(@stack) do
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

define provision_record(@raw_record) return @resource do
  $obj = to_object(@raw_record)
  $to_create = $obj["fields"]
  @resource = dme.record.create($to_create)
end

define delete_record(@record) do
  @record.destroy()
end