# design/media_types/do.rb
module V1
  module MediaTypes
    class DoCloud < Praxis::MediaType

      identifier 'application/json'

      attributes do
        attribute :id, Attributor::String
        attribute :href, Attributor::String
        attribute :name, Attributor::String
        attribute :memory, Attributor::Integer
        attribute :vcpus, Attributor::Integer
        attribute :disk, Attributor::Integer
        attribute :image, Attributor::Integer
        attribute :available, Attributor::Boolean
        attribute :region, Attributor::String
        attribute :size, Attributor::String
        attribute :status, Attributor::String
        attribute :deployment, Attributor::String
        attribute :server_template_href, Attributor::String
        attribute :api_host, Attributor::String
        attribute :cloud, Attributor::String
      end

      view :default do
        attribute :id
        attribute :href
        attribute :name
        attribute :memory
        attribute :vcpus
        attribute :disk
        attribute :image
        attribute :available
        attribute :region
        attribute :status
      end

      view :link do
        attribute :href
      end

    end
  end
end
