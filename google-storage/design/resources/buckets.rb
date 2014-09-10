  module ApiResources
    class Buckets
      include Praxis::ResourceDefinition

      media_type MediaTypes::Bucket

      routing do
        prefix '/acct/:acct/buckets'
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
          get '/:id', name: :bucket_href
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
          attribute :b, Attributor::Struct, required: true do
            attribute :name, String, required: true
            attribute :predefinedAcl, String
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
          attribute :id, String, required: true
        end
        response :no_content, media_type: nil
        response :bad_request, media_type: 'text/plain'
      end
    end
  end
