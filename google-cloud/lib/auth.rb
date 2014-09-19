# Google Ruby client docs: https://github.com/google/google-api-ruby-client

require 'google/api_client'
require 'google/api_client/client_secrets'
require 'google/api_client/auth/installed_app'

# We cache the API specs from the discovery service using this prefix in the filenames
CACHED_API_PREFIX = ".cache."

# Name and version of the Google API in the Google API discovery service
APIS = {
  'sql' => ['sqladmin', 'v1beta3'],
  'drive' => ['', ''],
}

# Oauth scopes required for each service
SCOPES = {
  'sql' => ['https://www.googleapis.com/auth/sqlservice.admin'],
  'drive' => [],
}

module GoogleCloud

  # Internal method to create a Google Auth client that is ready-to-go
  def self.auth_client
      client_secrets = Google::APIClient::ClientSecrets.load
      Signet::OAuth2::Client.new({
        application_name:     "RS Self-Service Google Cloud Proxy",
        application_version:  "0.0.1",
        authorization_uri:    'https://accounts.google.com/o/oauth2/auth',
        token_credential_uri: 'https://accounts.google.com/o/oauth2/token',
        client_id:            client_secrets.client_id,
        client_secret:        client_secrets.client_secret,
        scope:                SCOPES,
        redirect_uri:         "https://localhost:9292/auth/redirect",
      })
  end

  # Return the redirect URL pointing to Google where the user gets to accept the auth request
  # Returns the URL as a string
  def self.auth_redirect(acct, project)
    cli = auth_client
    cli.authorization_uri.to_s
  end

  def self.encode_creds(creds)
    Base64.encode64(Zlib.deflate(Yajl::Encoder.encode(creds)))
  end

  def self.decode_creds(str)
    begin
      Yajl::Parser.parse(Zlib.inflate(Base64.decode64(str)))
    rescue StandardError => e
      halt 400, "Cannot decode authentication credentials (#{e})"
    end
  end

  # Gets the refresh_token from the authorization code
  def self.get_creds(code)
    cli = auth_client
    cli.code = code
    puts "Fetching refresh/access tokens: #{cli.generate_access_token_request}"
    cli.fetch_access_token!
    encode_creds(access_token) if cli.access_token
  end

  # Return an authenticated API client based on auth creds passed as a hash
  # Returns the client
  def self.client(creds)
    # The auth stuff is copied from
    # https://github.com/google/google-api-ruby-client-samples/blob/master/drive/drive.rb
    client = Google::APIClient.new(
      application_name: "RS Self-Service Google Cloud Proxy",
      application_version: "0.0.1",
    )
    client.authorization = creds
    client.authorization.fetch_access_token!
    client
  end

  $google_api = {} # cached API definitions

  # Return a reference to the API definition, which can be used to construct calls
  def self.api(service)
    return $gogle_api[service] if $gogle_api.key?(service)

    # Load cached discovered API, if it exists. This prevents retrieving the
    # discovery document on every run, saving a round-trip to API servers.
    f = CACHED_API_PREFIX + service.to_s
    if File.exists? f
      File.open(f) do |file|
        $google_api[service] = Marshal.load(file)
      end
    else
      $google_api[service] = Google::APIClient.new.discovered_api('sqladmin', 'v1beta3')
      File.open(f, 'w') do |file|
        Marshal.dump($google_api[service], file)
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
