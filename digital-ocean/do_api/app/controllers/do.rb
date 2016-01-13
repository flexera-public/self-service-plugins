# app/controllers/do.rb
require 'logger'

module V1
  class DoCloud
    include Praxis::Controller
    implements V1::ApiResources::DoCloud

    DO_TOKEN = ENV['DO_TOKEN']
    RS_TOKEN = ENV['RS_REFRESH_TOKEN']
    DO_DROPLET_API = "https://api.digitalocean.com/v2/droplets"

    def create(**params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      do_userdata = "#!/bin/bash

sudo apt-get install curl
cd /tmp
wget https://rightlink.rightscale.com/rll/uca-0.3.0/rightlink.enable.sh
chmod +x ./rightlink.enable.sh
SERVER_NAME='" + request.payload.name + "'
SERVER_TEMPLATE='" + request.payload.server_template_href + "'
DEPLOYMENT_HREF='" + request.payload.deployment + "'         
CLOUD_NAME='" + request.payload.cloud + "'
HOST='" + request.payload.api_host + "'
KEY='" + RS_TOKEN + "'
sudo ./rightlink.enable.sh -l -n \\\"$SERVER_NAME\\\" -k \\\"$KEY\\\" -r \\\"$SERVER_TEMPLATE\\\" -c uca -a \\\"$HOST\\\" -e \\\"$DEPLOYMENT_HREF\\\" -f \\\"$CLOUD_NAME\\\"
"
      do_uri = DO_DROPLET_API
      do_request = '{"name":"' + request.payload.name + 
        '","region":"' + request.payload.region + 
        '","size":"' + request.payload.size + 
        '","image":' + request.payload.image.to_s + 
        ',"user_data":"' + do_userdata + '"}'

      do_response = do_post(do_uri, do_request)

      resp = JSON.parse(do_response.body)["droplet"]
      resp["href"] = "/api/do_proxy/droplets/" + resp["id"].to_s
      response = Praxis::Responses::Created.new()
      response.headers['Content-Type'] = 'application/json'
      response.headers['Location'] = resp["href"]
      response.body = resp
      response
    end

    def list(**other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      do_uri = DO_DROPLET_API
      do_response = do_get(do_uri)

      response.headers['Content-Type'] = 'vnd.rightscale.droplet+json;type=collection'

      resp = JSON.parse(do_response.body)["droplets"]
      resp.each do |r|
        r["href"] = "/api/do_proxy/droplets/" + r["id"].to_s
      end
      response.body = resp
      response
    end

    def show(id:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      do_uri = DO_DROPLET_API + "/" + id.to_s
      do_response = do_get(do_uri)

      resp = JSON.parse(do_response.body)["droplet"]
      resp["href"] = "/api/do_proxy/droplets/" + resp["id"].to_s
      response.headers['Content-Type'] = 'application/json'
      response.body = resp
      response
    end

    def powercycle(id:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      do_uri = DO_DROPLET_API + "/" + id.to_s + "/actions"
      do_request = '{"type":"power_cycle"}'

      do_response = do_post(do_uri, do_request)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def poweroff(id:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      do_uri = DO_DROPLET_API + "/" + id.to_s + "/actions"
      do_request = '{"type":"power_off"}'

      do_response = do_post(do_uri, do_request)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def delete(id:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      app = Praxis::Application.instance
      do_uri = DO_DROPLET_API + "/" + id.to_s
      app.logger.info("DELETE:" + do_uri)

      do_uri = URI.parse(do_uri)
      do_http = Net::HTTP.new(do_uri.host, do_uri.port)
      do_http.use_ssl = true
      ####do_http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      do_response = do_http.delete(do_uri,
         'Content-Type' => 'application/x-www-form-urlencoded',  
         'Authorization' => 'Bearer ' + DO_TOKEN)
      app.logger.info(do_response.body)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def do_post(do_uri, do_request)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      app = Praxis::Application.instance
      app.logger.info("GET:" + do_uri)
      do_uri = URI.parse(do_uri)
      do_http = Net::HTTP.new(do_uri.host, do_uri.port)
      do_http.use_ssl = true
      ####do_http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      do_response = do_http.post(do_uri, do_request, 
        'Content-Type'=>'application/json', 
        'Authorization' => 'Bearer ' + DO_TOKEN)
      app.logger.info(do_response.body)
      do_response
    end

    def do_get(do_uri)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      app = Praxis::Application.instance
      app.logger.info("POST:" + do_uri)

      do_uri = URI.parse(do_uri)
      do_http = Net::HTTP.new(do_uri.host, do_uri.port)
      do_http.use_ssl = true
      ####do_http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      do_response = do_http.get(do_uri, 
         'Authorization' => 'Bearer ' + DO_TOKEN)
      app.logger.info(do_response.body)
      do_response
    end

    private

    def authenticate!(secret)
      if secret != ENV["PLUGIN_SHARED_SECRET"]
        self.response = Praxis::Responses::Forbidden.new()
        response.body = { error: '403: Invalid shared secret'}
        return response
      else
        return nil
      end
    end

  end
end
