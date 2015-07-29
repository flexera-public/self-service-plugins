module V1
  module ApiResources
    class Change
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::Change
      version '1.0'
      prefix '/changes'

      action :show do
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
end
