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

ENV["RACK_ENV"] ||= "development"

require "./application.rb"

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

# recursively mount all controllers in a directory
def mount_dir(route_set, path, prefix)
  Dir.entries(path).each do |file|
    next if file[0] == '.' # skip silently
    if file =~ /\.rb$/
      require "#{path}/#{file}"
      resource_name = file.split(".").first
      r_class = resource_name.split('_').map{|e| e.capitalize}.join
      $logger.info "Loading resource type #{r_class} as #{prefix}#{resource_name}"
      route_set.add_route(Object.const_get(r_class),
          { :path_info => Rack::Mount::Strexp.compile("#{prefix}#{resource_name}",{},[],false) })
    elsif Dir.exist?("#{path}/#{file}")
      singular = file.sub(/s$/, "")
      mount_dir(route_set, "#{path}/#{file}", "#{prefix}#{file}/:#{singular}/")
    else
      $logger.warning "Skipping ./app/#{file}"
    end
  end
end


Routes = Rack::Mount::RouteSet.new do |set|
  mount_dir(set, "./app", "/")
end

run Routes
