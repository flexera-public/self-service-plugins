# See http://unicorn.bogomips.org/Unicorn/Configurator.html for complete
# documentation.

APP_PATH = File.expand_path(File.join(File.dirname(__FILE__), '..'))

if [ 'staging', 'production' ].include?(ENV['RACK_ENV'])
  PID_FILE = '/var/run/cisco-aci.pid'
  LOG_FILE = '/var/log/cisco-aci.log'
else
  PID_FILE = File.join(APP_PATH, 'rainbows.pid')
  LOG_FILE = APP_PATH + "/log/cisco-aci.log"
end

Rainbows! do
  use :ThreadPool # concurrency model to use
  worker_connections 5
  keepalive_timeout 2 # zero disables keepalives entirely
  client_max_body_size 5*1024*1024 # 5 megabytes

  #Rainbows.module_eval do
  #  EventMachine.kqueue = false if RUBY_PLATFORM =~ /darwin/i
  #end

end

# Help ensure your application will always spawn in the symlinked
# "current" directory that Capistrano sets up.
working_directory APP_PATH

worker_processes 2 # make this configurable later

# Create a nice logger and set the rainbows log destination to it. This will be fed through rack
# to sinatra to request.logger
$logger = ::Logger.new(LOG_FILE)
$logger.formatter = proc do |severity, datetime, progname, msg|
  "#{datetime}: #{msg}\n"
end
logger $logger
$logger.info("Rainbows starting")

# Redirect stdout and stderr to the same file as the logger
# In application.rb we then redirect these things to the logger itself
# If we do that here we get into trouble with forking workers...
stderr_path LOG_FILE # necessary to print exceptions which don't get redirected by the stuff below
stdout_path LOG_FILE # necessary to print exceptions which don't get redirected by the stuff below

# Requests can take a looonng time
timeout 120

# we use a shorter backlog for quicker failover when busy
listen "127.0.0.1:8001", :tcp_nopush => true

# PID file location
pid PID_FILE

# Don't preload app: minimal savings and breaks HUP process control signal
preload_app false

GC.respond_to?(:copy_on_write_friendly=) and GC.copy_on_write_friendly = true
