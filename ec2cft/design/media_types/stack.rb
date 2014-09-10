module V1
  module MediaTypes
    class Stack < Praxis::MediaType

      identifier 'application/json'

      attributes do
        attribute :id, Attributor::String, description: "Unique stack identifier"
        attribute :href, Attributor::String
        attribute :creation_time, Attributor::DateTime, description: "The time the stack was created"
        attribute :template_description, Attributor::String, description: "User defined description associated with the stack"
        attribute :name, Attributor::String, description: "Returns the stack name"
        attribute :status, Attributor::String, description: "The status of the stack"
        attribute :status_reason, Attributor::String, description: "Success/Failure message associated with the status"
        attribute :template, Attributor::String, description: "Returns the stack's template as a JSON string"
      end

      view :default do
        attribute :id
        attribute :href
        attribute :creation_time
        attribute :template_description
        attribute :name
        attribute :status
        attribute :status_reason
      end

      view :source do
        attribute :id
        attribute :href
        attribute :creation_time
        attribute :template_description
        attribute :name
        attribute :status
        attribute :status_reason
        attribute :template
      end

      view :link do
        attribute :href
      end
    end
  end
end
