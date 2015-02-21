module V1
  class Configuration
    include Praxis::Controller

    implements V1::ApiResources::Configuration

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
      href = "/api/accounts/#{account_id}/configuration/#{id}"
      details = { id: id,
                  bootstrap_script: script,
                  href: href }
      response.body = details
      response.headers['Location'] = href
      db[id] = details
      flush_db(db)
      response
    end

    def get_db
      if File.exists?('db.json')
        JSON.load(File.read('db.json'))
      else
        File.open('db.json', 'w') { |f| f.write({}.to_json) }
        {}
      end
    end

    def flush_db(db)
      File.open('db.json', 'w') { |f| f.write(db.to_json) }
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
      script = <<-SCRIPT
#!/bin/bash
# Download and install chef
curl -L https://www.chef.io/chef/install.sh | sudo bash
# Prepare the runlist file
cat <<'EOF' >/etc/chef/runlist.json
#{runlist_json}
EOF

mkdir -p /etc/chef

# Create the chef client configuration file
cat <<'EOF' >/etc/chef/client.rb
chef_server_url        "#{input[:chef_server_url]}"
validation_client_name "#{input[:validation_client_name]}"
node_name              "#{input[:node_name]}"
environment            "#{input[:chef_environment]}"
EOF

# Create the validation key file
cat <<'EOF' > /etc/chef/validation.pem
#{input[:validation_key]}
EOF

# Converge
chef-client -j /etc/chef/runlist.json
SCRIPT
      script
    end

  end
end
