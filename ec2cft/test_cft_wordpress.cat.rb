name "Test for CloudFormation Wordpress Template"
rs_ca_ver 20131202
short_description "Test for a CloudFormation Template

![logo](http://d2wwfe3odivqm9.cloudfront.net/wp-content/uploads/2014/05/aws_icon-cloudformation_white-200x200.png)"

# Cloud Selection
parameter "cft_url" do
  type "string"
  label "URL of the CloudFormation Template"
  default "https://s3-us-west-2.amazonaws.com/cloudformation-templates-us-west-2/WordPress_Single_Instance.template"
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
  template $cft_url
  parameters do {
    "KeyName" => $cft_param_key,
    "InstanceType" => $cft_param_instance_type,
    "DBRootPassword" => $cft_dbrootpw,
    "DBUser" => $cft_dbuser,
    "DBPassword" => $cft_dbpw
  } end
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