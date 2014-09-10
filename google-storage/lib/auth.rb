# Google Ruby client docs: https://github.com/google/google-api-ruby-client

require 'google/api_client'
require 'google/api_client/client_secrets'
require 'google/api_client/auth/installed_app'

CACHED_API_FILE = ".storage.cache"
CREDS_DIR       = ".gc_auth"
SCOPES          = ['https://www.googleapis.com/auth/devstorage.read_write']

# Include GoogleCloudStorageMixing into all controllers that require access to authenticated Google
# Cloud Storage APIs. This will add before-action hooks that ensure that the auth context is set,
# including @gc_storage_client and @gc_storage_api instance variables
module GoogleCloudStorageMixin
  def self.included(klass)
    klass.class_eval do
      # TODO: this before-action block should return a result for the errors instead of raising
      before :action do |controller|
        acct = controller.request.params.acct
        raise "Authentication is missing" unless acct
        client, project = GoogleCloudStorage.client(acct)
        raise "Authentication failed" unless client
        api = GoogleCloudStorage.api
        raise "Internal error: cannot retrieve Cloud Storage API definition" unless api
        puts "GoogleCloudStorageMixin: acct=#{acct} project=#{project}"
        controller.instance_variable_set(:@gc_storage_client, client)
        controller.instance_variable_set(:@gc_storage_project, project)
        controller.instance_variable_set(:@gc_storage_api, api)
        nil
      end
    end
  end
end

module GoogleCloudStorage

  $auth_cache = {} # hash of account id -> auth'd client

  # Internal method to create a Google Auth client that is ready-to-go
  def self.auth_client(acct, project)
      client_secrets = Google::APIClient::ClientSecrets.load
      Signet::OAuth2::Client.new({
        authorization_uri:    'https://accounts.google.com/o/oauth2/auth',
        token_credential_uri: 'https://accounts.google.com/o/oauth2/token',
        client_id:            client_secrets.client_id,
        client_secret:        client_secrets.client_secret,
        scope:                SCOPES,
        redirect_uri:         "https://localhost:9292/acct/#{acct}/auth/redirect?project=#{project}",
      })
  end

  # Verify that we have a valid authorization for the specified account, returns true/false
  def self.auth_test(acct)
    return true if $auth_cache.key?(acct)
    file_storage = FileStorage.new(CREDS_DIR, acct)
    return !file_storage.authorization.nil?
  end

  # Return the redirect URL pointing to Google where the user gets to accept the auth request
  # Returns the URL as a string
  def self.auth_redirect(acct, project)
    cli = auth_client(acct, project)
    cli.authorization_uri.to_s
  end

  # Saves the authorization code that Google responded with as a result of a successful auth
  # The project is simply saved with the auth code as a convenience
  # Returns true if the auth code was found to be valid
  def self.auth_save(acct, project, code)
    cli = auth_client(acct, project)
    cli.code = code
    #puts "Fetching access token: #{cli.generate_access_token_request}"
    cli.fetch_access_token!
    if cli.access_token
      file_storage = FileStorage.new(CREDS_DIR, acct)
      if file_storage.respond_to?(:write_credentials)
        file_storage.write_credentials(cli, project)
      end
      true
    else
      false
    end
  end

  # Return an authenticated API client based on auth creds stored in a file, i.e., the
  # whole auth thing is expected to have previously happened.
  # Returns the client and the project name
  def self.client(account)
    return $auth_cache[account] if $auth_cache.key?(account)

    # The auth stuff is copied from
    # https://github.com/google/google-api-ruby-client-samples/blob/master/drive/drive.rb
    client = Google::APIClient.new(
      application_name: "RS Self-Service Google Cloud Proxy",
      application_version: "0.0.1",
    )

    # FileStorage stores auth credentials in a file, so they survive multiple runs
    # of the application. This avoids prompting the user for authorization every
    # time the access token expires, by remembering the refresh token.
    creds_file = CREDS_DIR + "/#{account}"
    file_storage = FileStorage.new(CREDS_DIR, account)
    if file_storage.authorization.nil?
      [nil, nil]
    else
      client.authorization = file_storage.authorization
      client.authorization.fetch_access_token!
      $auth_cache[account] = [client, file_storage.project]
    end
  end

  $google_storage = nil # cached API definition

  # Return a reference to the API definition, which can be used ot construct calls
  def self.api
    return $gogle_storage if $gogle_storage

    # Load cached discovered API, if it exists. This prevents retrieving the
    # discovery document on every run, saving a round-trip to API servers.
    if File.exists? CACHED_API_FILE
      File.open(CACHED_API_FILE) do |file|
        $google_storage = Marshal.load(file)
      end
    else
      $google_storage = Google::APIClient.new.discovered_api('storage', 'v1')
      File.open(CACHED_API_FILE, 'w') do |file|
        Marshal.dump($google_storage, file)
      end
    end

    $google_storage
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
