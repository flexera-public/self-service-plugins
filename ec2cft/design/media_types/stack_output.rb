module V1
  module MediaTypes
    class StackOutput < Praxis::MediaType

      identifier 'application/json'

      attributes do
        attribute :key, Attributor::String, description: "Unique stack identifier"
        attribute :value, Attributor::String
        attribute :description, Attributor::String, description: "The time the stack was created"
      end

      view :default do
        attribute :key
        attribute :value
        attribute :description
      end

    end
  end
end
