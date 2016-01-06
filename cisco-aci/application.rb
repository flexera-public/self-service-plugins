#
# Copyright 2013 by RightScale, Inc. All Rights Reserved
#

require 'pry'

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
$apic_url = 'https://10.10.1.49'
$username = 'admin'
$password = 'rightscale11'

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

=begin
  before do
    #puts "request #{request.path} is: #{request.inspect}"
    #binding.pry
    creds_str = request.cookies['google-cloud']
    unless creds_str || (request.get? && request.path =~ /^\/auth/)
      $logger.info "google-cloud cookie missing, path is #{request.path_info}"
      halt 400, "google-cloud cookie missing"
    end

    if creds_str
      begin
        creds = GoogleCloud.decode_creds(creds_str)
        @client = GoogleCloud.client(creds)
      rescue StandardError => e
        halt 400, "Cannot decode authentication credentials (#{e})"
      end
    end

    # we allow a header to set the project
    request.params[:project] ||= request.env['X_PROJECT']
    $logger.info "Project: #{request.params[:project]}"
  end
=end

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

      def add_stuff(obj, stuff)
        #$logger.debug "Add stuff #{obj.class_name}: #{stuff.inspect}"
        stuff.each_pair do |k,v|
          # if it's a property, just set it
          if obj.props.key?(k)
            obj.set_prop(k, v)
            next
          end

          # see whether it's a child relationship resource
          cap_k = k[0].capitalize + k[1,k.length-1]
          cc = obj.child_classes.select{|cc| cc == cap_k || cc =~ /Rs#{cap_k}\z/} # also Rt?
          halt 400, "Ambiguous child class #{k} in #{obj.ruby_class}, choices: #{cc.sort.join(' ')}" \
            if cc.size > 1
          v.sub!(%r{^/mo.*/}, '') # convert value from href to name
          if cc.size == 1
            cc = cc.first
            child = Object.const_get("ACIrb::#{cc}").new(obj)
            if child.props.key?("name")
              child.set_prop('name', v)
            else
              name_props = child.props.select{|p,v| p.end_with?('Name')}
              if name_props.size == 1
                #$logger.debug "Setting prop #{name_props.first[0]}=#{v} (#{name_props.inspect})"
                child.set_prop(name_props.first[0], v)
              else
                halt 400, "Cannot set name for link '#{k}': #{name_props.sort.join(' ')}"
              end
            end
            next
          end

          halt 400, "Oops: #{obj.class_name} does not have attribute or child class #{k},\n" +
            "valid attributes: #{obj.props.keys.sort.join(' ')},\n" +
            "valid child classes: #{obj.child_classes.sort.join(' ')}"
        end
        obj
      end

      def gen_json(obj)
        if obj.is_a?(Array)
          obj.map{|o|o.to_json}
        else
          obj.to_json
        end
      end


  end

  # Helpers that can be called from outside a sinatra controller class: use Admin.method_name
  class << self

  end
end

#binding.pry
