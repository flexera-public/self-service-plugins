  module ApiResources
    class Auth
      include Praxis::ResourceDefinition

      media_type nil

      routing do
        prefix '/acct/:acct/auth'
      end

      action :show do
        use :has_account

        routing do
          get ''
        end

        params do
          attribute :project, String, required: true
        end

        response :ok, media_type: nil
        response :temporary_redirect, media_type: nil
        response :bad_request, media_type: nil
      end

      action :update do
        use :has_account

        routing do
          get '/redirect'
        end

        params do
          attribute :project, String, required: true
          attribute :code, String, required: true
        end

        response :ok, media_type: nil
        response :bad_request, media_type: nil
      end
    end
  end
