module SharedSecretAuthenticationPlugin
  include Praxis::PluginConcern

  class Plugin < Praxis::Plugin
    include Singleton

    def prepare_config!(node)
      node.attributes do
        attribute :authentication_default, Attributor::Boolean, default: true,
          description: 'Require authentication for all actions?'
      end
    end

    def config_key()
      :authentication
    end

    def authenticate(request)
      ENV['API_SHARED_SECRET'] && request.headers['X_Api_Shared_Secret'] == ENV['API_SHARED_SECRET']
    end
  end

  module Controller
    extend ActiveSupport::Concern

    included do
      before :action do |controller|
        action = controller.request.action
        if action.authentication_required
          unless Plugin.instance.authenticate(controller.request)
            Praxis::Responses::Unauthorized.new(body: 'unauthorized')
          end
        end
      end
    end
  end

  module ActionDefinition
    extend ActiveSupport::Concern

    included do
      decorate_docs do |action, docs|
        docs[:authentication_required] = action.authentication_required
      end
    end

    def requires_authentication(value)
      @authentication_required = value
    end

    def authentication_required
      @authentication_required ||= true
    end
  end
end
