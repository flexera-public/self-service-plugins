class Restifier < App
  $path = "/home/src/aws-sdk-core-ruby/aws-sdk-core/apis"
  $connector = "http://localhost:9001"

  $services = {}
  $client = nil

  APP_JSON = { 'content-type' => 'application/json' }

  helpers do

    def log_info
      $logger.info "Params   : #{params.inspect}"
      body_str = request.body.is_a?(String) ? request.body : request.body.read
      $logger.info "Body     : #{request.content_type} (#{body_str.size} bytes)"
      $logger.info "Service  : #{params['service']}"
      $logger.info "Res Type : #{params['resource_type']}"
    end

    def get_service(name)
      return $services[name] if $services.key?(name)
      svc_name = name.size < 4 ? name.upcase : name.camel_case
      $logger.info "Loading service #{svc_name}"
      a = Analyzer::Analyzer.new(path: $path, cloud: 'aws')
      svc = a.service_definition(svc_name)
      unless svc
        halt 404, "Service #{name} (#{svc_name}) is not supported"
      end
      $services[name] = svc
      $services[name]
    end

    def get_resource(svc, svc_name, name)
      #resource_name = Analyzer::AWS::ResourceRegistry.canonical_name(name)
      resource_name = name.singularize
      resource = svc.resources[resource_name]
      unless resource
        resources = svc.resources.keys.sort.join(' ')
        halt(400, "Service #{svc_name} does not have resource #{name}. " +
             "Available resources: #{resources}")
      end
      resource
    end

    def get_action(resource, res_name, name)
      action = resource.actions[name]
      unless action
        halt 400, "Resource #{res_name} does not have action #{name}"
      end
      $logger.info "Action: #{action.inspect}"
      action
    end

    def process_body(action, req_body)
      body_str = req_body.is_a?(String) ? req_body : req_body.read
      body = {}
      if body_str.size > 0
        unless request.content_type && request.content_type.start_with?("application/json")
          halt 400, "Request body content-type must be application/json"
        end
        begin
          body = Yajl::Parser.parse(request.body, :symbolize_keys => false)
        rescue StandardError => e
          halt 400, "Error parsing json body: #{e}"
        end
        unless body.is_a?(Hash)
          halt 400, "Request body must consist of a json hash"
        end
        $logger.info "Request body contains: #{body.keys.sort.join(' ')}"
      end
      body
    end

    def get_id(resource, id)
      return { resource.primary_id => id }
    end

    def perform_request(svc_name, action, req_body)
      url = URI("%s/%s/%s" % [$connector, svc_name, action.name])
      $logger.info "URL: #{url}"
      $logger.info "Payload: #{req_body.keys.join(' ')}"

      begin
        unless $client
          $client = Net::HTTP.start(url.host, url.port)
        end
        body_str = Yajl::Encoder.encode(req_body)
        body_str = "" if body_str == '{}'
        response = $client.request_post(url.path, body_str, APP_JSON)
      rescue Exception => e
        $logger.info "*** Error: #{e} #{e.inspect}"
        $logger.info e.backtrace[0..10].join(' | ')
        halt 500, e.message
      end

      $logger.info "Got: #{response.code} #{response.body}"
      return response.code.to_i, APP_JSON, response.body
    end

  end

  # show
  get '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "show")
    return perform_request(params['service'], action, {})
  end

  # index
  get '/:service/:resource_type' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "index")
    return perform_request(params['service'], action, {})
  end

  # custom collection actions
  post '/:service/:resource_type/action/:action' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], params['action'])
    body = process_body(action, request.body)
    return perform_request(params['service'], action, body)
  end

  # custom resource actions
  post '/:service/:resource_type/:id/action/:action' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], params['action'])
    body = process_body(action, request.body)
    body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, body)
  end

  # create
  post '/:service/:resource_type' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "create")
    body = process_body(action, request.body)
    return perform_request(params['service'], action, body)
  end

  # update
  put '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "update")
    body = process_body(action, request.body)
    body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, body)
  end

  # patch
  patch '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "patch")
    body = process_body(action, request.body)
    body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, body)
  end

  # delete
  delete '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "delete")
    body = process_body(action, request.body)
    body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, body)
  end

end
