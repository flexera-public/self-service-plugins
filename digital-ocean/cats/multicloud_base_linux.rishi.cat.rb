name "Base Linux for AWS/GCE/VMware"
rs_ca_ver 20131202
short_description "Allows you to launch a base Ubuntu or CentOS distro on AWS, GCE or VMware

![logo](http://icons.iconarchive.com/icons/carlosjj/google-jfk/128/compute-engine-icon.png) ![logo](http://www.sonru.com/images/uploads/global_aws.gif)

![logo](http://www.kronos.cn/uploadedImages/Global_Content/Partner_logos/VMwareLogoSmall.png)
"

long_description "This Cloud Application Template allows you to launch a base Ubuntu or CentOS distro on AWS, GCE or VMware. It uses the prescribed instance sizes (small, medium and large) for development purposes.

![logo](http://www.sonru.com/images/uploads/global_aws.gif)

![logo](http://icons.iconarchive.com/icons/carlosjj/google-jfk/128/compute-engine-icon.png)

![logo](http://www.kronos.cn/uploadedImages/Global_Content/Partner_logos/VMwareLogoSmall.png)

"

#########
# Parameters
#########

# Cloud
parameter "cloud" do
  type "string"
  label "Cloud"
  category "Infrastructure Providers"
  allowed_values "Amazon Web Services", "Google Compute Engine", "VMware"
  default "Google Compute Engine"
  description "Pick the cloud where the instance should launch."
end

# Instance Size
parameter "instance_size" do
  type "string"
  label "Instance Size"
  category "Infrastructure Providers"
  allowed_values "Small", "Medium", "Large"
  default "Small"
  description "Pick a reasonable instance size."
end

# Operating System
parameter "operating_system" do
  type "string"
  label "Distro"
  category "Operating System"
  allowed_values "Ubuntu", "CentOS"
  default "Ubuntu"
  description "Pick the distro."
end

#########
# Mappings
#########

# User-friendly cloud name to internal cloud name
# Also includes all the cloud specific configuration information
mapping "cloud_mapping" do {
  "Amazon Web Services" => {
    "cloud_href" => "/api/clouds/1",
    "cloud_name" => "EC2 us-east-1",
    "ssh_key" => "default",
    "datacenter" => "us-east-1e",
    "Ubuntu_mci" => "RightImage_Ubuntu_12.04_x64_v14.0.0",
    "Ubuntu_mci_rev" => "15",
    "CentOS_mci" => "RightImage_CentOS_6.5_x64_v14.0.0",
    "CentOS_mci_rev" => "14"
  },
  "Google Compute Engine" => {
    "cloud_name" => "Google",
    "cloud_href" => "/api/clouds/2175",
    "ssh_key" => null,
    "datacenter" => "us-central1-a",
    "Ubuntu_mci" => "RightImage_Ubuntu_12.04_x64_v14.0.0",
    "Ubuntu_mci_rev" => "15",
    "CentOS_mci" => "RightImage_CentOS_6.5_x64_v14.0.0",
    "CentOS_mci_rev" => "14"
  },
  "VMware" => {
    "cloud_name" => "VMware",
    "cloud_href" => "/api/clouds/2974",
    "ssh_key" => "rightscale_test",
    "datacenter" => "first-example-zone1",
    "Ubuntu_mci" => "RightImage_Ubuntu_12.04_x64_v14.0.0_vSphere",
    "Ubuntu_mci_rev" => "6",
    "CentOS_mci" => "RightImage_CentOS_6.5_x64_v14.0.0_vSphere",
    "CentOS_mci_rev" => "6"
  }
}
end

# User-friendly instance size names to appropriate sizes for the clouds
mapping "instance_size_mapping" do {
  "Small" => { 
    "Amazon Web Services" => "c3.large",
    "Google Compute Engine" => "n1-standard-1",
    "VMware" => "small"
  },
  "Medium" => { 
    "Amazon Web Services" => "c2.xlarge",
    "Google Compute Engine" => "n1-standard-2",
    "VMware" => "medium"
  },
  "Large" => { 
    "Amazon Web Services" => "c3.2xlarge",
    "Google Compute Engine" => "n1-standard-4",
    "VMware" => "large"
  }
}
end

#########
# Resources
#########

resource "base_server", type: "server" do
  name                  "Base Linux Server"
  cloud_href            map($cloud_mapping, $cloud, "cloud_href")
  server_template       find("Base ServerTemplate for Linux (v14.0.1)", revision: 39)
  multi_cloud_image     find(map( $cloud_mapping, $cloud, join([$operating_system, "_mci"])), 
                            revision: map($cloud_mapping, $cloud, join([$operating_system, "_mci_rev"])))
  instance_type         find(map($instance_size_mapping, $instance_size, $cloud), cloud_href: map($cloud_mapping, $cloud, "cloud_href"))
  ssh_key               find(resource_uid: map($cloud_mapping, $cloud, "ssh_key"), cloud_href: map($cloud_mapping, $cloud, "cloud_href"))
  datacenter            find(map($cloud_mapping, $cloud, "datacenter"), cloud_href: map($cloud_mapping, $cloud, "cloud_href"))
end


