module V1
  module MediaTypes
    class Record < Praxis::MediaType

      identifier 'application/vnd.rightscale.record+json'
      @@kind = 'route53#record'

      attributes do
        attribute :kind, String
        attribute :id, String # Maybe not so much?
        attribute :href, String
        attribute :name, String
        attribute :type, String
        attribute :values, Attributor::Collection.of(String)
        attribute :change, Change

        links do
          link :change
        end
      end

      view :default do
        attribute :kind
        attribute :id
        attribute :name
        attribute :type
        attribute :values
        attribute :links
      end

      view :link do
        attribute :href
      end

      def href()
        V1::ApiResources::Record.prefix+'/'+id
      end

      def kind()
        @@kind
      end

    end

    class RecordCollectionSummary < Praxis::MediaType
      attributes do
        attribute :href, String
      end

      view :default do
        attribute :href
      end

      view :link do
        attribute :href
      end
    end
  end
end
