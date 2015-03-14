module V1
  module MediaTypes
    class Configuration < Praxis::MediaType
      identifier 'application/json'

      attributes do
        attribute :id, String,
          description: 'The ID of this configuration'
        attribute :href, String,
          description: 'The HREF of this configuration'
        attribute :bootstrap_script, String,
          description: 'The bootstrap script that can be used for configuration'
      end

      view :default do
        attribute :id
        attribute :href
        attribute :bootstrap_script
      end

      view :link do
        attribute :href
      end
    end
  end
end
