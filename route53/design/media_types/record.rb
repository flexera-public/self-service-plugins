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
        attribute :ttl, Integer
        attribute :values, Attributor::Collection.of(String)
        attribute :change, Change
        attribute :links, Attributor::Collection.of(Hash)
      end

      view :default do
        attribute :kind
        attribute :id
        attribute :name
        attribute :href
        attribute :type
        attribute :ttl
        attribute :values
        attribute :links
      end

      view :link do
        attribute :href
      end

    end
  end
end
