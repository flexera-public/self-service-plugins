module V1
  module ApiResources
    class PublicZone
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::PublicZone
      version '1.0'
      prefix '/public_zones'
      trait :authorized

      action :index do
        routing do
          get ''
        end
        response :ok
      end

      action :show do
        routing do
          get '/:id'
        end
        params do
          attribute :id, String, required: true
        end
        response :ok
        response :not_found
        response :bad_request
      end

      action :create do
        routing do
          post ''
        end
        payload do
          attribute :name, String, required: true
        end
        response :created
        response :bad_request
      end

      action :delete do
        routing do
          delete '/:id'
        end
        params do
          attribute :id, String, required: true
        end
        response :no_content
        response :bad_request
        response :not_found
      end

      action :release do
        routing do
          post '/:id/release'
        end
        payload do
          attribute :name, String, required: true
        end
        response :ok
        response :bad_request
        response :not_found
      end

    end
  end
end
