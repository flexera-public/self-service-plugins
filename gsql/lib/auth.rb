# Google Ruby client docs: https://github.com/google/google-api-ruby-client

require 'google/api_client'
require 'google/api_client/client_secrets'
require 'google/api_client/auth/file_storage'
require 'google/api_client/auth/installed_app'

CACHED_API_FILE = ".sqladmin.cache"
CREDS_DIR       = ".gc_auth"

# Include GoogleCloudSQLMixing into all controllers that require access to authenticated Google
# Cloud SQL APIs. This will add before-action hooks that ensure that the auth context is set,
# including @gc_sql_client and @gc_sql_api
module GoogleCloudSQLMixin
  def self.included(klass)
    klass.class_eval do
      before :action do |controller|
        puts "GoogleCloudSQLMixin before action"
        acct = controller.request.params.acct
        raise "Authentication is missing" unless acct
        @gc_sql_client = GoogleCloudSQL.client(acct)
        raise "Authentication failed" unless @gc_sql_client
        @gc_sql_api = GoogleCloudSQL.api
        raise "Internal error: cannot retrieve Cloud SQL API definition" unless @gc_sql_api
      end
    end
  end
end

module GoogleCloudSQL

  $auth_cache = {} # hash of account id -> auth'd client

  def self.auth_client(acct)
      client_secrets = Google::APIClient::ClientSecrets.load
      Signet::OAuth2::Client.new({
        authorization_uri:    'https://accounts.google.com/o/oauth2/auth',
        token_credential_uri: 'https://accounts.google.com/o/oauth2/token',
        client_id:            client_secrets.client_id,
        client_secret:        client_secrets.client_secret,
        scope:                ['https://www.googleapis.com/auth/sqlservice.admin'],
        redirect_uri:         "https://localhost:9/acct/#{acct}/auth/redirect",
      })
  end

  def self.auth_test(acct)
    return true if $auth_cache.key?(acct)
    creds_file = CREDS_DIR + "/#{acct}"
    file_storage = Google::APIClient::FileStorage.new(creds_file)
    return !file_storage.authorization.nil?
  end

  def self.auth_redirect(acct)
    cli = auth_client(acct)
    cli.authorization_uri.to_s
  end

  def self.auth_set(acct, code)
    cli = auth_client(acct)
    cli.code = code
    puts "Fetching access token: #{cli.generate_access_token_request}"
    cli.fetch_access_token!
    if cli.access_token
      file_storage = Google::APIClient::FileStorage.new(creds_file)
      if file_storage.respond_to?(:write_credentials)
        file_storage.write_credentials(cli)
      end
      true
    else
      false
    end
  end

  # Handle authentication and load the API
  def self.client(account)
    return $auth_cache[account] if $auth_cache.key?(account)

    # The auth stuff is copied from
    # https://github.com/google/google-api-ruby-client-samples/blob/master/drive/drive.rb
    client = Google::APIClient.new(
      application_name: "RS Self-Service Google Cloud Proxy",
      application_version: "0.0.1",
    )

    unless File.directory?(CREDS_DIR)
      Dir.mkdir(CREDS_DIR)
    end

    # FileStorage stores auth credentials in a file, so they survive multiple runs
    # of the application. This avoids prompting the user for authorization every
    # time the access token expires, by remembering the refresh token.
    creds_file = CREDS_DIR + "/#{account}"
    file_storage = Google::APIClient::FileStorage.new(creds_file)
    if file_storage.authorization.nil?
      nil
    else
      client.authorization = file_storage.authorization
      $auth_cache[account] = client
      client
    end
  end

  $google_sql = nil # cached API definition

  def self.api
    return $gogle_sql if $gogle_sql

    # Load cached discovered API, if it exists. This prevents retrieving the
    # discovery document on every run, saving a round-trip to API servers.
    if File.exists? CACHED_API_FILE
      File.open(CACHED_API_FILE) do |file|
        $google_sql = Marshal.load(file)
      end
    else
      $google_sql = client.discovered_api('sqladmin', 'v1beta3')
      File.open(CACHED_API_FILE, 'w') do |file|
        Marshal.dump($google_sql, file)
      end
    end

    $google_sql
  end

end

# This is no longer used...
=begin
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

=end
