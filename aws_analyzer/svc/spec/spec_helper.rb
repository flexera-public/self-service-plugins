require 'sinatra'
require 'rack/test'
require 'logger'
require 'sinatra/json'
require 'yajl'
require 'multi_json'
require 'base64'
require "net/http"
require "net/https"
require "uri"
require_relative '../../analyzer'
require_relative "../application.rb"
require_relative "../app/restifier"

ENV["RACK_ENV"] ||= "development"

$logger = ::Logger.new(STDERR)
$logger.info "Hello Logger"

RSpec.configure do |config|
  config.include Rack::Test::Methods

  def app
    Restifier
  end
end

def post_json(uri, args)
  post uri, Yajl::Encoder.encode(args), "CONTENT_TYPE" => "application/json"
end

def put_response(resp)
  if resp.status == 201
    puts "OK, Location: #{resp.location}"
  elsif resp.status < 300
    if resp.body.size > 0
      puts "OK: #{JSON.pretty_generate(JSON.parse(resp.body))}"
    else
      puts "OK (empty body)"
    end
  else
    puts "ERROR #{resp.status}: #{resp.body}"
  end
end

