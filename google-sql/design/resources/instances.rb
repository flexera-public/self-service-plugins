  module ApiResources
    class Instances
      include Praxis::ResourceDefinition

      media_type MediaTypes::Instance

      routing do
        prefix '/acct/:acct/instances'
      end

      action :index do
        use :has_account

        routing do
          get ''
        end
        response :ok
        response :bad_request, media_type: 'text/plain'
      end

      action :show do
        use :has_account

        routing do
          get '/:id', name: :instance_href
        end
        params do
          attribute :id, String, required: true
        end
        response :ok
        response :bad_request, media_type: 'text/plain'
      end

      action :create do
        use :has_account

        routing do
          post ''
        end
        params do
        end
        payload do
          attribute :instance, String, required: true
          attribute :settings, Attributor::Struct, required: true do
            attribute :tier, String, required: true
          end
        end
        response :created, media_type: nil
        response :bad_request, media_type: 'text/plain'
      end

      action :delete do
        use :has_account

        routing do
          delete '/:id'
        end
        params do
        end
        response :no_content, media_type: nil
        response :bad_request, media_type: 'text/plain'
      end
    end
  end
