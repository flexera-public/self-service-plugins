  module ApiResources
    class Auth
      include Praxis::ResourceDefinition

      media_type MediaTypes::Auth

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

        response :ok
        response :temporary_redirect
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

        response :ok
      end
    end
  end
