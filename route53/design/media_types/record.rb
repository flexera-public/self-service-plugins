module V1
  module MediaTypes
    class Record < Praxis::MediaType

      identifier 'application/vnd.rightscale.record+json'

      attributes do
        attribute :kind, String
        attribute :id, String
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
        attribute :href
        attribute :type
        attribute :values
        attribute :links
      end

      view :link do
        attribute :href
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
