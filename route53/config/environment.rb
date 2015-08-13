# Main entry point - DO NOT MODIFY THIS FILE
ENV['RACK_ENV'] ||= 'development'

Bundler.require(:default, ENV['RACK_ENV'])

require File.expand_path('../../app/plugins/shared_secret_authentication_plugin', __FILE__)

# Default application layout.
# NOTE: This layout need NOT be specified explicitly.
# It is provided just for illustration.
Praxis::Application.instance.layout do
  map :initializers, 'config/initializers/**/*'
  map :lib, 'lib/**/*'
  map :design, 'design/' do
    map :api, 'api.rb'
    map :media_types, '**/media_types/**/*'
    map :resources, '**/resources/**/*'
  end
  map :app, 'app/' do
    map :plugins, 'plugins/**/*'
    map :models, 'models/**/*'
    map :controllers, '**/controllers/**/*'
    map :responses, '**/responses/**/*'
    map :helpers, '**/helpers/**/*'
    map :attributes, '**/attributes/**/*'
  end
end

Praxis::Application.configure do |application|
  application.bootloader.use SharedSecretAuthenticationPlugin
end
