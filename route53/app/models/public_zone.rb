module V1
  module Models
    class PublicZone

      attr_accessor :kind, :id, :href, :change

      attr_reader :hosted_zone

      def initialize(hosted_zone, change_info=nil)
        @hosted_zone = hosted_zone
        @kind = 'route53#public_zone'
        @id = hosted_zone.id.match(/\/[a-z_]*\/([a-z0-9A-Z_]*)$/)[1]
        @href = V1::ApiResources::PublicZone.prefix+'/'+@id
        @href = '/'+ENV['SUB_PATH']+@href if ENV.has_key?('SUB_PATH')
        # @name = hosted_zone.name
        # @caller_reference = hosted_zone.caller_reference
        # @config = hosted_zone.config
        # @resource_record_set_count = hosted_zone.resource_record_set_count
        @change = change_info if change_info
      end

      def records_summary()
        OpenStruct.new(href: href+'/records')
      end

      def links()
        links = []
        links << { rel: 'self', href: href }
        links << { rel: 'records', href: href+'/records' }

        if @change
          change_href = V1::ApiResources::Change.prefix+'/'+@change['id']
          change_href = '/'+ENV['SUB_PATH']+change_href if ENV.has_key?('SUB_PATH')
          links << { rel: 'change', href: change_href}
        end
      end

      def method_missing(m, **args, &block)
        @hosted_zone.send(m)
      end
    end
  end
end
