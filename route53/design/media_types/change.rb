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
        attribute :links, Attributor::Collection.of(Hash)
      end

      view :default do
        attribute :kind
        attribute :id
        attribute :href
        attribute :status
        attribute :submitted_at
        attribute :links
      end

      view :link do
        attribute :href
      end

      def href()
        href = V1::ApiResources::Change.prefix+'/'+id
        href = '/'+ENV['SUB_PATH']+href if ENV.has_key?('SUB_PATH')
      end

      def links()
        links = []
        links << { rel: 'self', href: href }
        links
      end

      def kind()
        @@kind
      end
    end
  end
end
