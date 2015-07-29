module V1
  module Models
    class PublicZone

      attr_accessor :kind, :id, :href, :name,
        :caller_reference, :config, :resource_record_set_count, :change

      def initialize(hosted_zone, change_info=nil)
        @kind = 'route53#public_zone'
        @id = hosted_zone.id.match(/\/[a-z_]*\/([a-z0-9A-Z_]*)$/)[1]
        @href = V1::ApiResources::PublicZone.prefix+'/'+@id
        @name = hosted_zone.name
        @caller_reference = hosted_zone.caller_reference
        @config = hosted_zone.config
        @resource_record_set_count = hosted_zone.resource_record_set_count
        @change = change_info if change_info
      end

      def records_summary()
        OpenStruct.new(href: href+'/records')
      end
    end
  end
end
