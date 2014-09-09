#! /home/tve/.rbenv/shims/ruby
# This command line utility was used to prototype some of the Google Auth stuff
# It's probably not useful anymore.

# Google Ruby client docs: https://github.com/google/google-api-ruby-client

require 'google/api_client'
require 'google/api_client/client_secrets'
require 'google/api_client/auth/file_storage'
require 'google/api_client/auth/installed_app'

CACHED_API_FILE = ".sqladmin.cache"
CREDS_FILE      = ".gc_auth"

module Google
  class APIClient
    class InstalledAppFlow
      def authorize_cli(storage)
        puts "Please visit: #{@authorization.authorization_uri.to_s}"
        printf "Enter the code: code="
        code = gets
        @authorization.code = code
        @authorization.fetch_access_token!
        if @authorization.access_token
          if storage.respond_to?(:write_credentials)
            storage.write_credentials(@authorization)
          end
          @authorization
        else
          nil
        end
      end
    end
  end
end

# Handle authentication and load the API
def setup
  # The auth stuff is copied from
  # https://github.com/google/google-api-ruby-client-samples/blob/master/drive/drive.rb

  client = Google::APIClient.new(
    application_name: "RS Self-Service Google Cloud Proxy",
    application_version: "0.0.1",
  )

  # FileStorage stores auth credentials in a file, so they survive multiple runs
  # of the application. This avoids prompting the user for authorization every
  # time the access token expires, by remembering the refresh token.
  # Note: FileStorage is not suitable for multi-user applications.
  file_storage = Google::APIClient::FileStorage.new(CREDS_FILE)
  if file_storage.authorization.nil?
    client_secrets = Google::APIClient::ClientSecrets.load
    # The InstalledAppFlow is a helper class to handle the OAuth 2.0 installed
    # application flow, which ties in with FileStorage to store credentials
    # between runs.
    flow = Google::APIClient::InstalledAppFlow.new(
      :client_id => client_secrets.client_id,
      :client_secret => client_secrets.client_secret,
      :scope => ['https://www.googleapis.com/auth/sqlservice.admin']
    )
    client.authorization = flow.authorize_cli(file_storage)
  else
    client.authorization = file_storage.authorization
  end

  sql = nil
  # Load cached discovered API, if it exists. This prevents retrieving the
  # discovery document on every run, saving a round-trip to API servers.
  if File.exists? CACHED_API_FILE
    File.open(CACHED_API_FILE) do |file|
      sql = Marshal.load(file)
    end
  else
    sql = client.discovered_api('sqladmin', 'v1beta3')
    File.open(CACHED_API_FILE, 'w') do |file|
      Marshal.dump(sql, file)
    end
  end

  return client, sql
end

def list_instances(client, sqladmin)
  result = client.execute(
    api_method: sqladmin.instances.list,
    parameters: { project: 'rightscale.com:tve-test' },
  )
  [result.status, result.data.to_json]
end

$client, $sqladmin = setup
puts list_instances($client, $sqladmin)




