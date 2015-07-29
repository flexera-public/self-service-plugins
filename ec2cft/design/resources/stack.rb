module V1
  module ApiResources
    class Stack
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::Stack
      version '1.0'
      prefix '/ec2cft/accounts/:account_id/stacks'
      trait :authenticated
      
      action :index do
        use :versionable

        routing do
          get ''
        end
        params do
          attribute :account_id, Attributor::Integer, required: true
        end
        response :ok
        response :forbidden
      end

      action :show do
        use :versionable

        routing do
          get '/:name'
        end
        params do
          attribute :account_id, Attributor::Integer, required: true
          attribute :name, required: true
          attribute :view, Attributor::String, default: "default"
        end
        response :ok
        response :not_found
        response :forbidden
      end

      action :create do
        routing do
          post ''
        end        

        params do
          attribute :account_id, Attributor::Integer, required: true
        end

        payload do
          attribute :name, required: true
          attribute :template, Attributor::String, required: true
          attribute :parameters, Attributor::Hash
        end

        response :created
        response :unprocessable_entity
        response :forbidden
      end

      action :delete do
        routing do
          delete '/:name'
        end

        params do
          attribute :account_id, Attributor::Integer, required: true
          attribute :name, required: true
        end

        response :no_content
        response :unprocessable_entity
        response :forbidden
      end

    end


  end
end


