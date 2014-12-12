require 'rubygems'
require 'bundler/setup'
require 'yaml'
require 'active_support/inflector'
Dir["#{File.dirname(__FILE__)}/lib/**/*.rb"].each { |f| load(f) }
