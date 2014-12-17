#
# Copyright 2013 by RightScale, Inc. All Rights Reserved
#

require 'pry'

class String
  def camel_case
    return self if self !~ /_/ && self =~ /[A-Z]+.*/
    split('_').map{|e| e.capitalize}.join
  end
end

class App < Sinatra::Base
  helpers Sinatra::JSON

  error 500 do
    if env['sinatra.error']
      $logger.info "***** BOOM *****"
      $logger.info env['sinatra.error'].message
      $logger.info "\n" + env['sinatra.error'].backtrace[0..20].join("\n")
    end
  end

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
    $logger.info("Processing #{self.class}# #{request.env["REQUEST_URI"].inspect} " +
        "(for #{request.env["HTTP_X_FORWARDED_FOR"]}) [#{request.env["REQUEST_METHOD"]}] ")
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
