# -p 8080

require 'bundler/setup'
require 'praxis'
require 'thin'

application = Praxis::Application.instance
application.logger = Logger.new(STDOUT)
application.setup

run application
