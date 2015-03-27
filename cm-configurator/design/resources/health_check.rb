module V0
  module ApiResources
    class HealthCheck
      include Praxis::ResourceDefinition

      media_type 'text/plain'

      routing do
        prefix "/health-check"
      end

      action :health_check do
        routing do
          get ""
        end

        response :ok
      end
    end
  end
end
