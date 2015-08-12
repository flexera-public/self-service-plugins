name "Namespace and sanity test of route53 plugin"
rs_ca_ver 20131202
short_description "Namespance and sanity test of route53 plugin"

resource "dns_zone", type: "route53.public_zone" do
  name "foo.bar.com"
end

resource "dns_record", type: "route53.record" do
  zone @dns_zone
  name "yeah"
  type "A"
  ttl 60
  values ["1.2.3.4","5.6.7.8"]
end

# operation "launch" do
#   describe "Handle the dependencies between public_zone and record"
#   definition "launch"
# end
#
# define launch(@dns_zone, @dns_record) do
# end

# Creates a "log" entry in the form of an audit entry.  The target of the audit
# entry defaults to the deployment created by the CloudApp, but can be specified
# with the "auditee_href" option.
#
# @param $summary [String] the value to write in the "summary" field of an audit entry
# @param $options [Hash] a hash of options where the possible keys are;
#   * detail [String] the message to write to the "detail" field of the audit entry. Default: ""
#   * notify [String] the event notification catgory, one of (None|Notification|Security|Error).  Default: None
#   * auditee_href [String] the auditee_href (target) for the audit entry. Default: @@deployment.href
#
# @see http://reference.rightscale.com/api1.5/resources/ResourceAuditEntries.html#create
define sys_log($summary,$options) do
  $log_default_options = {
    detail: "",
    notify: "None",
    auditee_href: @@deployment.href
  }

  $log_merged_options = $options + $log_default_options
  rs.audit_entries.create(
    notify: $log_merged_options["notify"],
    audit_entry: {
      auditee_href: $log_merged_options["auditee_href"],
      summary: $summary,
      detail: $log_merged_options["detail"]
    }
  )
end

namespace "route53" do
  service do
    host "https://route53plugin.cse.rightscale-services.com" # HTTP endpoint presenting an API defined by self-serviceto act on resources
    path "/route53"  # path prefix for all resources, RightScale account_id substituted in for multi-tenancy
    headers do {
      "user-agent" => "self_service" ,     # special headers as needed
      "X-Api-Version" => "1.0",
      "X-Api-Shared-Secret" => "7aed5e1f-f8f8-40e3-81ee-15a8a8868242"
    } end
  end
  type "public_zone" do
    fields do
      field "name" do
        type "string"
        required true
      end
    end
  end

  type "record" do
    provision "provision_record"
    fields do
      field "zone" do
        type "resource"
      end
      field "name" do
        type "string"
        required true
      end
      field "type" do
        type "string"
        required true
      end
      field "ttl" do
        type "number"
        required true
      end
      field "values" do
        type "array"
        required true
      end
    end
  end
end

define provision_record(@raw_record) return @resource do
  call sys_log("Raw Record", {detail: to_object(@raw_record.zone)})
end
