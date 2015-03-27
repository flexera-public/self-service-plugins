# See http://unicorn.bogomips.org/Unicorn/Configurator.html for complete
# documentation.

require 'bundler/setup'

APP_PATH = File.expand_path(File.join(File.dirname(__FILE__), '..'))

nr_processes = 8
nr_connections = 25

Rainbows! do
  if nr_connections > 1
    use :ThreadPool
  end

  worker_processes nr_processes
  worker_connections nr_connections
  keepalive_timeout 5 # zero disables keepalives entirely
  client_max_body_size 5*1024*1024 # 5 megabytes
end

# All production STs include an application owner field (which we assume to also be the same group)
# Thus use that user and group accordingly to run the application.
if [ 'integration', 'staging', 'production', 'meta' ].include?(ENV['RACK_ENV'])
  user('www-data','www-data')
end

# Help ensure your application will always spawn in the symlinked
# "current" directory that Capistrano sets up.
working_directory APP_PATH

# listen on both a Unix domain socket and a TCP port,
# we use a shorter backlog for quicker failover when busy
# listen 8080, :tcp_nopush => true
listen 8080

# Enable preloading for copy-on-write memory savings.
# Do this only if debugging is disabled.
# preload_app !debug

GC.respond_to?(:copy_on_write_friendly=) and GC.copy_on_write_friendly = true

after_fork do |server, worker|
  # re-establish db connections here
end
