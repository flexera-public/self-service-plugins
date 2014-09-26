# app/controllers/do.rb
require 'logger'

module V1
  class DoCloud
    include Praxis::Controller
    implements V1::ApiResources::DoCloud

    DO_TOKEN = ENV['TOKEN']
    DO_DROPLET_API = "https://api.digitalocean.com/v2/droplets"

    def create(**params)
      do_uri = DO_DROPLET_API
      do_request = '{"name":"' + request.payload.name + 
        '","region":"' + request.payload.region + 
        '","size":"' + request.payload.size + 
        '","image":' + request.payload.image.to_s + '}'

      do_response = do_post(do_uri, do_request)

      resp = JSON.parse(do_response.body)["droplet"]
      resp["href"] = "/api/do_proxy/droplets/" + resp["id"].to_s
      response.headers['Content-Type'] = 'application/json'
      response.headers['Location'] = resp["href"]
      response.body = resp
      response
    end

    def list(**other_params)
      do_uri = DO_DROPLET_API
      do_response = do_get(do_uri)

      response.headers['Content-Type'] = 'vnd.rightscale.droplet+json;type=collection'

      ##### binding.pry
      resp = JSON.parse(do_response.body)["droplets"]
      resp.each do |r|
        r["href"] = "/api/do_proxy/droplets/" + r["id"].to_s
      end
      response.body = resp
      response
    end

    def show(id:, **other_params)
      do_uri = DO_DROPLET_API + "/" + id.to_s
      do_response = do_get(do_uri)

      resp = JSON.parse(do_response.body)["droplet"]
      resp["href"] = "/api/do_proxy/droplets/" + resp["id"].to_s
      response.headers['Content-Type'] = 'application/json'
      response.body = resp
      response
    end

    def powercycle(id:, **other_params)
      do_uri = DO_DROPLET_API + "/" + id.to_s + "/actions"
      do_request = '{"type":"power_cycle"}'

      do_response = do_post(do_uri, do_request)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def poweroff(id:, **other_params)
      do_uri = DO_DROPLET_API + "/" + id.to_s + "/actions"
      do_request = '{"type":"power_off"}'

      do_response = do_post(do_uri, do_request)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def delete(id:, **other_params)
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

  end
end
