# design/media_types/do.rb
module V1
  module MediaTypes
    class DoCloud < Praxis::MediaType

      identifier 'application/json'

      attributes do
        attribute :id, String
        attribute :name, String
        attribute :memory, Integer
        attribute :vcpus, Integer
        attribute :disk, Integer
        attribute :image, Integer
        attribute :available, Attributor::Boolean
        attribute :region, String
        attribute :size, String
        attribute :status, String
      end

      view :default do
        attribute :id
        attribute :name
        attribute :memory
        attribute :vcpus
        attribute :disk
        attribute :image
        attribute :available
        attribute :region
        attribute :status
      end
    end
  end

end
