name "Test for DME IP address"
rs_ca_ver 20131202
short_description "Test for a DME A name

![logo](http://www.dnsmadeeasy.com/wp-content/uploads/2013/09/logo1.png)"


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
parameter "domain_name" do
  type "string"
  label "Name for A record"
  category "DNS"
end

# Cloud Selection
parameter "domain_ip" do
  type "string"
  label "IP address for A record"
  category "DNS"
end

resource 'dns_entry', type: 'dme.record' do
  domain 'dev.rightscaleit.com' # alternatively: domain_id 1234565
  name $domain_name
  type "A"
  value $domain_ip
  ttl 30
end

operation 'Update IP Address' do
  definition 'update_ip'
  description 'Update the IP of the record'
end

define update_ip(@dns_entry, $domain_ip) return @dns_entry do
 @dns_entry.update(name: @dns_entry.name, type: @dns_entry.type, value: $domain_ip)
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
    host "http://54.227.94.207:8080"        # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/dme/accounts/:account_id"      # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service" ,     # special headers as needed
      "X-Api-Version" => "1.0"
    } end
  end
  type "record" do                          
    provision "provision_record"            
    delete "delete_record"                  
    # path "/records" # Unneeded since we'll use the name of the type by default
    fields do
      domain do
        type "string"
        required true
      end
      name do
        type "string"
        required true
      end
      value do
        type "string"
        required true
      end
      type do
        type "string"
        required true
      end
      dynamicDns do
        type "boolean"
      end
      ttl do
        type "number"
      end
    end
    # outputs ["domain", "name", "value", "type", "ttl"]
  end
end 

define provision_record(@raw_record) return @resource do
  @resource = dme.record.create(record: to_object(@raw_record))
end

define delete_record(@record) do
  @record.destroy()
end