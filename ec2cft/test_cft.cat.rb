name "Test for CloudFormation Template"
rs_ca_ver 20131202
short_description "Test for a CloudFormation Template

![logo](http://d2wwfe3odivqm9.cloudfront.net/wp-content/uploads/2014/05/aws_icon-cloudformation_white-200x200.png)"


# output 'dns_name' do
#   label "Domain name"
#   category "General"
#   default_value join([@dns_entry.name,".dev.rightscaleit.com"])
# end

# output 'dns_ip' do
#   label "IP Address"
#   category "General"
#   default_value @dns_entry.value
# end

# Cloud Selection
parameter "cft_url" do
  type "string"
  label "URL of the CloudFormation Template"
  default "https://s3-external-1.amazonaws.com/cloudformation-samples-us-east-1/Rails_Simple.template"
end

resource 'cft_template', type: 'ec2cft.stack' do
  name @@deployment.name
  template $cft_url
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
    host "http://54.227.94.207:80"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/ec2cft/accounts/:account_id"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service" ,     # special headers as needed
      "X-API-Version" => "1.0"
    } end
  end
  type "stack" do                          
    provision "provision_record"            
    delete "delete_record"                  
    # path "/records" # Unneeded since we'll use the name of the type by default
    fields do
      template do
        type "string"
        required true
      end
      name do
        type "string"
        required true
      end
    end
    # outputs ["domain", "name", "value", "type", "ttl"]
  end
end 

define provision_record(@raw_stack) return @resource do
  @resource = ec2cft.stack.create(record: to_object(@raw_stack))
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