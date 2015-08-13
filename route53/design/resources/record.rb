module V1
  module ApiResources
    class Record
      include Praxis::ResourceDefinition

      media_type V1::MediaTypes::Record
      version '1.0'
      trait :authorized

      action_defaults do
        routing do
          prefix '//public_zones/:public_zone_id/records'
        end
        params do
          attribute :public_zone_id, String
        end
      end

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
          post '//records'
          post ''
        end
        payload required: true do
          attribute :public_zone_id, String
          attribute :name, String, required: true
          attribute :type, String,
            required: true,
            values: [
              "SOA",
              "A",
              "TXT",
              "NS",
              "CNAME",
              "MX",
              "PTR",
              "SRV",
              "SPF",
              "AAAA"
            ]
          attribute :ttl, Integer, required: true
          attribute :values, Attributor::Collection.of(String), required: true
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
      end

      action :release do
        routing do
          post '/:id/release'
        end
        params do
          attribute :id, String, required: true
        end
        response :ok
        response :bad_request
      end

    end
  end
end
