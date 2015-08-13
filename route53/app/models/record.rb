require 'digest/sha1'

module V1
  module Models
    class Record

      attr_accessor :kind, :href, :change

      attr_reader :record

      def initialize(zone_id, record, change_info=nil)
        @record = record
        @kind = 'route53#record'
        @zone_href = V1::ApiResources::PublicZone.prefix+'/'+zone_id
        @href = @zone_href+
          V1::ApiResources::Record.prefix+'/'+id
        @href = '/'+ENV['SUB_PATH']+@href if ENV.has_key?('SUB_PATH')
        @change = change_info if change_info
      end

      def method_missing(m, **args, &block)
        @record.send(m)
      end

      def values
        @record.resource_records.map{ |r| r.value }
      end

      def id
        to_sha = "#{@record.name}#{@record.type}#{@record.values.first}"
        Digest::SHA1.hexdigest(to_sha)
      end

      def links()
        links = []
        links << { rel: 'self', href: href }
        links << { rel: 'public_zone', href: @zone_href }
        links
      end
    end
  end
end
