require 'rest-client'

module V1
  class ChefConfiguration
    include Praxis::Controller

    implements V1::ApiResources::ChefConfiguration

    def show(account_id:, id:, **other_params)
      db = get_db
      if (conf = db[id])
        response.body = conf
      else
        self.response = Praxis::Responses::NotFound.new()
        response.body = { error: 'Could not find configuration with given id' }
      end
      response.headers['Content-Type'] = 'application/json'
      response
    end

    def create(account_id:, **other_params)
      db = get_db
      script = generate(request.payload.to_h)
      self.response = Praxis::Responses::Created.new()
      id = BSON::ObjectId.new.to_s
      href = "/chef_configurations/#{id}"
      details = { id: id,
                  kind: 'cm-configuration#chef',
                  bootstrap_script: script,
                  href: href }
      response.body = details
      response.headers['Location'] = href
      db[id] = details
      flush_db(db)
      response
    end

    # Get the database
    def get_db
      if File.exists?('/var/db') && (content = JSON.load(File.read('/var/db'))) == nil
        File.open('/var/db', 'w') { |f| f.write({}.to_json) }
        {}
      else
        content
      end
    end

    def flush_db(db)
      File.open('/var/db', 'w') { |f| f.write(db.to_json) }
    end

    # Generates the cloud-init script for installing and configuring the server using Chef
    #
    # @param input [Hash] options for generating the cloud-init script
    #
    # @option options [String] :chef_server_url
    # @option options [String] :validation_client_name
    # @option options [String] :validation_key
    # @option options [Array<String>] :run_list
    # @option options [String] :chef_environment
    # @option options [Hash] :attributes
    #
    def generate(input)
      runlist_json = JSON.pretty_generate((input[:attributes] || {}).merge({ run_list: input[:run_list] }))
      validation_key = get_cred(input[:validation_key])
      script = <<-SCRIPT
#!/bin/bash
# Download and install chef
curl -L https://www.chef.io/chef/install.sh | sudo bash
mkdir -p /etc/chef

# Prepare the runlist file
cat <<'EOF' >/etc/chef/runlist.json
#{runlist_json}
EOF

# Create the chef client configuration file
cat <<'EOF' >/etc/chef/client.rb
chef_server_url        "#{input[:chef_server_url]}"
validation_client_name "#{input[:validation_client_name]}"
node_name              "#{input[:node_name]}"
environment            "#{input[:chef_environment]}"
EOF

# Create the validation key file
cat <<'EOF' > /etc/chef/validation.pem
#{validation_key}
EOF

# Converge
chef-client -j /etc/chef/runlist.json
SCRIPT
      script
    end

    def get_cred(name)
      response = RestClient::Request.execute(
        method: :get,
        headers: {
          'X-Api-Version' => '1.5',
          'Authorization' => "Bearer #{ENV["RS_AUTH_TOKEN"]}",
          'Accept' => 'application/json',
          'params' => { 'filter[]' => "name==#{name}" }
        },
        url: "#{ENV["RS_API_ENDPOINT"]}/api/credentials")
      return nil unless response.code == 200
      response = JSON.parse(response)
      cred = response.detect { |c| c['name'] == name }
      return nil unless cred

      cred_href = (self_link = cred['links'].detect { |l| l['href'] if l['rel'] == 'self' }) && self_link['href']

      response = RestClient::Request.execute(
        method: :get,
        headers: {
          'X-Api-Version' => '1.5',
          'Authorization' => "Bearer #{ENV["RS_AUTH_TOKEN"]}",
          'Accept' => 'application/json',
          'params' => { 'view' => 'sensitive' }
        },
        url: "#{ENV["RS_API_ENDPOINT"]}#{cred_href}")
      return nil unless response.code == 200
      response = JSON.parse(response)
      response && response['value']
    end


  end
end
