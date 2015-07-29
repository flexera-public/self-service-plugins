module V1
  module MediaTypes
    class Change < Praxis::MediaType

      identifier 'application/vnd.rightscale.change+json'

      attributes do
        attribute :id, Attributes::Route53Id
        attribute :href, String
        attribute :status, String
        attribute :submitted_at, String
      end

      view :default do
        attribute :id
        attribute :href
        attribute :status
        attribute :submitted_at
      end

      view :link do
        attribute :href
      end

      def href()
        V1::ApiResources::Change.prefix+'/'+id
      end
    end
  end
end
