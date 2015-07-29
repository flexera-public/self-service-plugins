module V1
  module MediaTypes
    class Change < Praxis::MediaType

      identifier 'application/vnd.rightscale.change+json'
      @@kind = 'route53#change'

      attributes do
        attribute :kind, String
        attribute :id, Attributes::Route53Id
        attribute :href, String
        attribute :status, String
        attribute :submitted_at, String
      end

      view :default do
        attribute :kind
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

      def kind()
        @@kind
      end
    end
  end
end
