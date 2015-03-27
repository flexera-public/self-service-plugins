module V1
  module MediaTypes
    class ChefConfiguration < Praxis::MediaType
      identifier 'application/json'

      attributes do
        attribute :id, String,
          description: 'The ID of this configuration',
          example: '54fa382f636c6f5406b70100'
        attribute :kind, String,
          description: 'The kind of this configuration',
          example: 'cm-configuration#chef'
        attribute :href, String,
          description: 'The HREF of this configuration',
          example: '/api/accounts/60073/configurations/54fa382f636c6f5406b70100'
        attribute :bootstrap_script, String,
          description: 'The bootstrap script that can be used for configuration',
          example: ''
      end

      view :default do
        attribute :id
        attribute :kind
        attribute :href
        attribute :bootstrap_script
      end

      view :link do
        attribute :href
      end
    end
  end
end
