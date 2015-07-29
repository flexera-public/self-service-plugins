module V1
  module ApiResources
    class Record
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::Record
      version '1.0'
      prefix '/records'

      parent V1::ApiResources::PublicZone

      action :index do
        routing do
          get ''
        end
        response :ok
      end
    end
  end
end
