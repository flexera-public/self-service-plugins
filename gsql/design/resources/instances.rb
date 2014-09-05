  module ApiResources
    class Instances
      include Praxis::ResourceDefinition

      media_type MediaTypes::Instance

      routing do
        prefix '/acct/:acct'
      end

      action :index do
        use :has_account

        routing do
          get ''
        end
        response :ok
      end

      action :show do
        use :has_account

        routing do
          get '/:id'
        end
        params do
          attribute :id, String, required: true
        end
        response :ok
      end
    end
  end
