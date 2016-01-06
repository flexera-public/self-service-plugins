# design/resources/do.rb
module V1
  module ApiResources
    class DoCloud
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::DoCloud
      version "1.0"
      trait :authenticated

      routing do
        prefix "/api/do_proxy/droplets"
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
          attribute :deployment
          attribute :server_template_href
          attribute :api_host
          attribute :cloud
        end

        response :created
      end # create

      action :list do
        routing do
          get ''
        end

        response :ok, media_type: "vnd.rightscale.droplet"
      end #show

      action :powercycle do
        routing do
          post "/:id/actions/powercycle"
        end

        params do
          attribute :id
        end 
        response :ok, media_type: "application/json"
      end #powercycle

      action :poweroff do
        routing do
          post "/:id/actions/poweroff"
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
