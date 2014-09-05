module V1
  module ApiResources
    class Hello
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::Hello
      version '1.0'

      routing do
        prefix '/api/hello'
      end

      action :index do
        use :versionable

        routing do
          get ''
        end
        response :ok
      end

      action :show do
        use :versionable

        routing do
          get '/:id'
        end
        params do
          attribute :id, Integer, required: true, min: 0
        end
        response :ok
      end
    end
  end
end
