#\ -p 8888

require 'bundler/setup'
require 'praxis'
require File.expand_path('../app/plugins/shared_secret_authentication_plugin', __FILE__)

application = Praxis::Application.instance
application.logger = Logger.new(STDOUT)
application.setup

run application
