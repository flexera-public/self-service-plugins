  module ApiResources
    class Instances
      include Praxis::ResourceDefinition

      media_type MediaTypes::Instance

      routing do
        prefix '/acct/:acct/instances'
      end

      action :index do
        use :has_account
        use :versionable

        routing do
          get ''
        end
        response :ok
        response :bad_request, media_type: 'text/plain'
      end

      action :show do
        use :has_account
        use :versionable

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
        use :versionable

        routing do
          post ''
        end

        payload do
          #attribute :i, Attributor::Struct, required: true do
            attribute :instance, String, required: true
            attribute :masterInstanceName, String
            attribute :region, String
            attribute :tier, String, required: true # should be in settings...
            attribute :settings, Attributor::Struct, required: true do
              attribute :activationPolicy, String
              attribute :authorizedGaeApplications, Attributor::Collection.of(String)
              attribute :backupConfiguration, Attributor::Collection.of(Attributor::Struct) do
                attribute :binaryLogEnabled, Attributor::Boolean
                attribute :enabled, Attributor::Boolean
                attribute :startTime, String
              end
              attribute :databaseFlags, Attributor::Collection.of(String)
              attribute :ipConfiguration, Attributor::Collection.of(Attributor::Struct) do
                attribute :authorizedNetworks, Attributor::Collection.of(String)
                attribute :enabled, Attributor::Boolean
                attribute :requireSsl, Attributor::Boolean
              end
              attribute :locationPreference, Attributor::Collection.of(Attributor::Struct) do
                attribute :followGaeApplication, String
                attribute :zone, String
              end
              attribute :pricingPlan, String
              attribute :replicationType, String
            end
          end
        #end
        response :created, media_type: nil
        response :bad_request, media_type: 'text/plain'
      end

      action :delete do
        use :has_account
        use :versionable

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
