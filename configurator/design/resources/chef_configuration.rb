module V1
  module ApiResources
    class ChefConfiguration
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::ChefConfiguration
      version '1.0'

      routing do
        prefix '/cm/accounts/:account_id/chef_configurations'
      end

      action :show do
        use :versionable

        routing do
          get '/:id'
        end

        params do
          attribute :account_id, Integer, required: true
          attribute :id, String, required: true
        end

        response :ok
        response :not_found
      end

      action :create do
        routing do
          post ''
        end

        params do
          attribute :account_id, String, required: true
        end

        payload do
          attribute :chef_server_url, String, required: true,
            description: 'The URL of the chef server',
            example: 'https://api.opscode.com/organizations/rs-st-dev'
          attribute :node_name, String, required: true,
            description: 'The name of the node',
            example: 'mysql_server'
          attribute :validation_client_name, String, required: true,
            description: 'The name of the validator',
            example: 'rs-st-dev-validator'
          attribute :validation_key, String, required: true,
            description: 'The validation key'
          attribute :chef_environment, String, required: true,
            description: 'The chef environment',
            example: 'staging'
          attribute :run_list, Attributor::Collection.of(String), required: true,
            description: 'The runlist for the chef client'
          attribute :first_attributes, Attributor::Hash.of(key: String, value: Object),
            description: 'Optional attributes for the chef run'
        end

        response :created
        response :unprocessable_entity
      end

    end
  end
end
