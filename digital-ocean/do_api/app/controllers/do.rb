# app/controllers/do.rb
module V1
  class DoCloud
    include Praxis::Controller

    # TODO:~ definitely delete this!
    DO_TOKEN = ""
    DO_DROPLET_API = "https://api.digitalocean.com/v2/droplets"

    implements V1::ApiResources::DoCloud

    def create(**params)
      do_uri = DO_DROPLET_API
      do_request = '{"name":"' + request.payload.name + 
        '","region":"' + request.payload.region + 
        '","size":"' + request.payload.size + 
        '","image":' + request.payload.image.to_s + '}'

      do_response = do_post(do_uri, do_request)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def list(**other_params)
      do_uri = DO_DROPLET_API
      do_response = do_get(do_uri)
      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def show(id:, **other_params)
      do_uri = DO_DROPLET_API + "/" + id.to_s
      do_response = do_get(do_uri)
      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
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
      do_uri = DO_DROPLET_API + "/" + id.to_s

      do_uri = URI.parse(do_uri)
      do_http = Net::HTTP.new(do_uri.host, do_uri.port)
      do_http.use_ssl = true
      ####do_http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      do_response = do_http.delete(do_uri,
         'Content-Type' => 'application/x-www-form-urlencoded',  
         'Authorization' => 'Bearer ' + DO_TOKEN)

      response.headers['Content-Type'] = 'application/json'
      response.body = do_response.body
      response
    end

    def do_post(do_uri, do_request)
      do_uri = URI.parse(do_uri)
      do_http = Net::HTTP.new(do_uri.host, do_uri.port)
      do_http.use_ssl = true
      ####do_http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      do_response = do_http.post(do_uri, do_request, 
        'Content-Type'=>'application/json', 
        'Authorization' => 'Bearer ' + DO_TOKEN)
      do_response
    end

    def do_get(do_uri)
      do_uri = URI.parse(do_uri)
      do_http = Net::HTTP.new(do_uri.host, do_uri.port)
      do_http.use_ssl = true
      ####do_http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      do_response = do_http.get(do_uri, 
         'Authorization' => 'Bearer ' + DO_TOKEN)
      do_response
    end

  end
end
