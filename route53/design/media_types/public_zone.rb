module V1
  module MediaTypes
    class PublicZone < Praxis::MediaType

      identifier 'application/vnd.rightscale.public_zone+json'
      @@kind = 'route53#public_zone'

      attributes do
        attribute :kind, String
        attribute :id, String
        attribute :href, String
        attribute :name, String
        attribute :caller_reference, String
        attribute :config do
          attribute :comment, String
          attribute :private_zone, String
        end
        attribute :resource_record_set_count, Integer

        links do
          link :change, Change
          link :records, V1::MediaTypes::RecordCollectionSummary, using: :records_summary
        end
      end

      view :default do
        attribute :kind
        attribute :id
        attribute :href
        attribute :name
        attribute :caller_reference
        attribute :config
        attribute :resource_record_set_count
        attribute :links
      end

      view :link do
        attribute :href
      end

      # def href()
      #   V1::ApiResources::PublicZone.prefix+'/'+id
      # end
      #
      # def kind()
      #   @@kind
      # end

      # def records()
      #   { href: href+'/records' }
      # end
    end
  end
end
