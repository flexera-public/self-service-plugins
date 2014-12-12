require 'rack/mount'
require 'logger'
require 'sinatra'
require 'sinatra/json'
require 'yajl'
require 'multi_json'
require 'base64'
require "net/http"
require "net/https"
require "uri"
require_relative '../analyzer'
require_relative "application.rb"

ENV["RACK_ENV"] ||= "development"

# load lib directory
Dir["./lib/*.rb"].each do |file|
  require file
end

# I do not believe that this 'takes' for Sinatra::Base stuff
#configure do
#  disable :show_exceptions
#end

# To support rackup (instead of rainbows)
$logger ||= ::Logger.new(STDERR)
#puts "config.ru"
#STDERR.puts "config.ru"

require_relative "app/restifier"
Routes = Rack::Mount::RouteSet.new do |set|
  set.add_route(Restifier, { path_info: %r{^/} }, {}, :restifier)
end

run Routes
