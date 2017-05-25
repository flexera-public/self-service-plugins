#
# Copyright 2013 by RightScale, Inc. All Rights Reserved
#

require 'pry'
require 'sinatra'
require 'sinatra/json'

# Redirect stderr/stdout through the logger
# This breaks forking, we get a IOError - uninitialized stream
=begin
class IOToLog < IO
  def initialize(logger); @logger = logger; end
  def write(string)
    @logger.info(string) unless string == "\n"
  end
end
$stderr = STDERR = IOToLog.new($logger)
$stdout = STDOUT = IOToLog.new($logger)
#puts "Hello stdout"
#STDERR.puts "Hello STDERR"
#$stderr.puts "Hello $stderr"
=end
#$apic_url = 'https://10.10.1.49'
$apic_url = ENV['APIC_URL']
$username = ENV['APIC_USER']
$password = ENV['APIC_PASS']

require 'acirb'
$api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
                             format: "json", debug: false)

class App < Sinatra::Base
  helpers Sinatra::JSON

  configure do
    disable :show_exceptions
    set :logging, nil  # this prevents sinatra from mucking with env['rack.logger']
    enable :logging   # this is useful when using rackup
  end

  # ensure we return plain text errors and not something that pretends to be text/html
  error 400..599 do
    content_type "plain/text"
    nil  # apparently this preserves the error returned by a 'halt' statement
  end

  before do
    halt 403, "Permission denied" unless request.env["HTTP_X_TOKEN"] == "cisco-aci-demo-20160110"
  end

  before do
    $request_start_time = Time.now
    # Logger troubleshooting
    #puts "ENV=#{request.env.keys.sort.join(' ')}"
    #puts "RACK_LOGGER=#{request.env['rack.logger'].class}"
    #puts "REQUEST_LOGGER=#{request.logger.inspect}"
    #request.logger.info("Processing #{self.class}# #{request.env["REQUEST_URI"].inspect} " +
    #    "(for #{request.env["HTTP_X_FORWARDED_FOR"]}) [#{request.env["REQUEST_METHOD"]}] ")
    $logger.info("***************");
    $logger.info("Processing #{self.class}# #{request.env["REQUEST_METHOD"]} #{request.env["REQUEST_URI"].inspect} " +
        "(for #{request.env["HTTP_X_FORWARDED_FOR"]})] ")
  end

  # for PUT/POST requests, if there is a json body then parse the params from the
  # json body instead of expecting that the query string will do
  before do
    params.merge!(env["rack.routing_args"]) if env["rack.routing_args"]
    if (request.request_method == "POST" || request.request_method == "PUT") &&
        request.content_type && request.content_type.start_with?("application/json")
      #$logger.info "JSON PUT/POST"
      begin
        parsed = request.body ? Yajl::Parser.parse(request.body) : {}
        #$logger.info "Body parsed: #{parsed.inspect}"
        if parsed.is_a?(Hash)
          $logger.info "Reading json content-body, merging #{parsed.keys.sort.join(' ')}"
          parsed.each_pair{|k,v| params[k.to_sym] = v} # .merge! doesn't work...
          $logger.info "Params after merge: #{params.keys.sort.join(' ')}"
        end
      rescue StandardError => e
        halt 400, "Error parsing json body: #{e}"
      end
    end
  end

  before do
    logger.info("Params: #{params.map{|k,v| "#{k}=\"#{v}\""}.join(", ")}") unless params.size == 0
  end

  after do
    $request_start_time ||= Time.now  # just in case the before filter didn't run
    $logger.info("Completed in #{Time.now - $request_start_time}s "+
        "| #{response.status} [#{request.env["REQUEST_URI"].inspect}]")
  end

  # Global helpers that can be called from within any controller class
  helpers do

  end

  # Helpers that can be called from outside a sinatra controller class: use Admin.method_name
  class << self

  end
end

#binding.pry
