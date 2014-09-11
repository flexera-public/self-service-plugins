name "Test for DME IP address and DigitalOcean instance"
rs_ca_ver 20131202
short_description "Test for a DME A name on a DigitalOcean instance

![logo](http://curvve-curvvemedia.netdna-ssl.com/wp-content/uploads/2013/08/Digital-Ocean-logo-tall-175x100.png) ![logo](http://www.dnsmadeeasy.com/wp-content/uploads/2013/09/logo1.png)
"


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

resource 'dns_entry', type: 'dme.record' do
  domain 'dev.rightscaleit.com' # alternatively: domain_id 1234565
  name $domain_name
  type "A"
  value "1.1.1.1"
  ttl 30
end

resource "do_instance", type: "do.droplet" do
  name                  @@deployment.name
  size                  "512mb"
  region                "sfo1"
  image                 5141286
end

operation "launch" do
  definition "launch"
  description "Launch the app"
end

define launch(@dns_entry, @do_instance) return @dns_entry, @do_instance do
  provision(@do_instance)
  provision(@dns_entry)
  call update_ip(@dns_entry, @do_instance.networks["v4"][0]["ip_address"]) retrieve @dns_entry
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
      "user-agent" => "self_service",      # special headers as needed
      "X-API-Version" => "1.0"
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


############################
############################
 #  ___  _      _ _        _    ___                   
 # |   \(_)__ _(_) |_ __ _| |  / _ \ __ ___ __ _ _ _  
 # | |) | / _` | |  _/ _` | | | (_) / _/ -_) _` | ' \ 
 # |___/|_\__, |_|\__\__,_|_|  \___/\__\___\__,_|_||_|
 #        |___/                                       
############################
############################

namespace "do" do
  service do
    host "ec2-54-202-222-194.us-west-2.compute.amazonaws.com"        # HTTP endpoint presenting an API defined by self-service to act on resources
    path "/api/do_proxy"                                             # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "X-Api-Version" => "1.0"                                       # special headers as needed
    } end
  end
  type "droplet" do                       # defines resource of type "droplet"
    provision "provision_droplet"         # name of RCL definition to use to provision the resource
    delete "delete_droplet"               # name of RCL definition to use to delete the resource
    fields do                             # field of a droplet with rules for validation
      name do                               
        type "string"
        required true
      end
      region do                               
        type "string"
        required true
      end
      size do                               
        type "string"
        required true
      end
      image do                               
        type "number"
        required true
      end
    end
  end
end

# Define the RCL definitions to create and destroy the resource
define provision_droplet(@raw_droplet) return @droplet do
  @droplet = do.droplet.create(droplet: to_object(@raw_droplet)) # Calls .create on the API resource
end
define delete_droplet(@droplet) do
  @droplet.destroy() # Calls .delete on the API resource
end