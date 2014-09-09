# design/resources/do.rb
module V1
  module ApiResources
    class DoCloud
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::DoCloud
      version "1.0"

      routing do
        prefix "/api/do_proxy/droplet"
      end

      action :create do
        routing do
          post ''
        end
        payload do
          attribute :name
          attribute :size
          attribute :image
          attribute :region
        end
        response :ok, media_type: "application/json"
      end # create

      action :list do
        routing do
          get ''
        end

        response :ok, media_type: "application/json"
      end #show

      action :powercycle do
        routing do
          get "/:id/powercycle"
        end

        params do
          attribute :id
        end 
        response :ok, media_type: "application/json"
      end #powercycle

      action :poweroff do
        routing do
          get "/:id/poweroff"
        end

        params do
          attribute :id
        end 
        response :ok, media_type: "application/json"
      end #poweroff

      action :delete do
        routing do
          delete "/:id"
        end

        params do
          attribute :id
        end 
        response :ok, media_type: "application/json"
      end #poweroff

      action :show do
        routing do
          get "/:id"
        end

        params do
          attribute :id
        end 
        response :ok, media_type: "application/json"
      end #show

    end
  end
end
