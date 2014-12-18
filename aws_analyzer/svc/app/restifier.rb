class Restifier < App
  $paths = ["/home/src/aws-sdk-core-ruby/aws-sdk-core/apis", "./apis" ]
  $connector = "http://localhost:9001"

  $services = {}
  $client = nil

  APP_JSON = { 'content-type' => 'application/json' }

  before do
    @body = parse_body(request.body, request.content_type)
    $logger.info "Request body contains: #{@body.keys.sort.join(' ')}"
  end

  helpers do

    def log_info
      $logger.info "Params   : #{params.inspect}"
      $logger.info "Service  : #{params['service']}"
      $logger.info "Res Type : #{params['resource_type']}"
    end

    def parse_body(body, content_type, name="request")
      body_str = body.is_a?(String) ? body : body.read
      $logger.info "Body     : #{content_type} (#{body_str.size} bytes)"

      body = {}
      if body_str.size > 0
        unless content_type && content_type.start_with?("application/json")
          halt 400, "#{name} body content-type must be application/json not #{content_type}"
        end
        begin
          body = Yajl::Parser.parse(body_str, :symbolize_keys => false)
        rescue StandardError => e
          halt 400, "Error parsing json body: #{e}"
        end
        unless body.is_a?(Hash)
          $logger.info "Body     : #{body}"
          halt 400, "#{name} body must consist of a json hash, not #{body.class}"
        end
      end
      body
    end

    def get_service(name)
      return $services[name] if $services.key?(name)
      svc_name = name.size < 4 ? name.upcase : name.camel_case
      $logger.info "Loading service #{svc_name}"
      a = Analyzer::Analyzer.new(paths: $paths, cloud: 'aws')
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

    def check_action(resource, res_name, name)
      action = resource.actions[name]
      $logger.info "Action   : #{action.name}" if action
      action
    end

    def get_action(resource, res_name, name)
      action = resource.actions[name] || resource.custom_actions[name]
      unless action
        al = (resource.actions.keys + resource.custom_actions.keys).sort.join(' ')
        halt 400, "Resource #{res_name} does not have action #{name}, available actions: #{al}"
      end
      $logger.info "Action   : #{action.name}"
      action
    end

    # Produce a has that has the resource's ID field set, this is used to construct the
    # resource's JSON when receiving the ID in the query string
    def get_id(resource, id)
      halt 500, "Resource metadata for #{resource.name} doesn't specify primary key" \
          unless resource.primary_id && resource.primary_id.size > 0
      return { resource.primary_id => id }
    end

    # Extract the primary key from a result payload
    def extract_id(resource, res)
      halt 500, "Resource metadata for #{resource.name} doesn't specify primary key" \
          unless resource.primary_id && resource.primary_id.size > 0
      halt 400, "Response with #{resource.name} does not have primary key #{resource.primary_id}" \
          unless res.key?(resource.primary_id)
      res[resource.primary_id]
    end

    def make_href(params, id)
      "/#{params[:service]}/#{params[:resource_type]}/#{id}"
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
      if response.code.to_i == 200 && (response.body.size == 0 || response.body == '{}')
        return 204, APP_JSON, nil
      else
        return response.code.to_i, APP_JSON, response.body
      end
    end

  end

  # show
  get '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = check_action(resource, params['resource_type'], "show")
    if action
      filter = { resource.primary_id => params[:id] }
    else
      action = check_action(resource, params['resource_type'], "index")
      halt 400, "Resource #{params['resource_type']} has neither show nor index action" unless action
      shape = svc.shapes[action.payload]
      $logger.info "Shape: #{shape.inspect}"
      if shape['type'] == 'structure' && shape['members'].key?(resource.primary_id)
        filter = { resource.primary_id => params[:id] }
      else
        filter = { resource.primary_id.pluralize => [ params[:id] ] }
      end
    end
    halt 500, "Resource #{params['resource_type']} does not have primary_id" unless resource.primary_id
    code, hdrs, body =  perform_request(params['service'], action, filter)
    if code == 200 && hdrs['content-type'] == "application/json"
      res = parse_body(body, hdrs['content-type'], "response")
      if res.key?(params['resource_type']) && res[params['resource_type']].is_a?(Array)
        $logger.info "Got hash with #{params['resource_type']} array"
        res = res[params['resource_type']]
      elsif res.key?(params['resource_type'].singularize)
        res = res[params['resource_type'].singularize]
      end
      if res.is_a?(Array)
        if res.size == 1
          res = res.first
        else
          halt 500, "Got response with #{res.size} elements: #{res.inspect}"
        end
      end
      if res.is_a?(Hash)
        if res.key?(resource.primary_id)
          res['links'] ||= []
          res['links'] << { rel: 'self', href: make_href(params, res[resource.primary_id]) }
          return 200, { "Content-Type" => "application/json" }, Yajl::Encoder.encode(res)
        else
          halt 500, "Response doesn't have primary key #{resource.primary_id}: #{res.inspect}"
        end
      end
      halt 500, "Can't interpret response: #{res.inspect}"
    else
      return code, hdrs, body # uhh, I'm sure we need to do something smart here
    end
  end

  # index
  get '/:service/:resource_type' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    halt 500, "Resource #{params['resource_type']} has no primary key" unless resource.primary_id
    action = get_action(resource, params['resource_type'], "index")
    args = {}
    if params['filter'] && params['filter'].is_a?(Array)
      $logger.info "Got filter"
      params['filter'].each do |f|
        k, v = f.split('==', 2)
      $logger.info "filter: #{k} == #{v}"
        args[k] = v
      end
    end
    code, hdrs, body =  perform_request(params['service'], action, args)
    if code == 200 && hdrs['content-type'] == "application/json"
      res = parse_body(body, hdrs['content-type'], "response")
      if res.is_a?(Hash)
        if res.key?(params['resource_type']) && res[params['resource_type']].is_a?(Array)
          res = res[params['resource_type']]
          $logger.info "Got hash with #{params['resource_type']} array"
        else
          arrays = res.values.select{|v| v.is_a?(Array)}
          if arrays.size != 1
            halt 500, "Can't interpret response: #{res.inspect}"
          end
        end
      end
      if res.is_a?(Array)
        res.map!{ |r|
          if r.is_a?(Hash) && r.key?(resource.primary_id)
            r['links'] ||= []
            r['links'] << { rel: 'self', href: make_href(params, r[resource.primary_id]) }
          else
            $logger.info "Oops? Missing #{resource.primary_id}"
          end # should error out if there's no primary key?
          r
        }
        return 200, { "Content-Type" => "application/json" }, Yajl::Encoder.encode(res)
      end
      halt 500, "Can't interpret response: #{res.inspect}"
    else
      return code, hdrs, body # uhh, I'm sure we need to do something smart here
    end
  end

  # custom service (top-level) action
  post '/:service/actions/:action' do
    log_info
    svc = get_service(params['service'])
    action = svc.actions[params['action']]
    $logger.info "Action is #{action.inspect}" if action
    $logger.info "Actions are #{svc.actions.inspect}" unless action
    halt(400, "Service #{params['service']} does not have an action #{params['action']}, " +
      "available actions: #{svc.actions.keys.join(' ')}") unless action
    $logger.info "Action   : #{action.name}" if action
    return perform_request(params['service'], action, @body)
  end

  # custom collection actions
  post '/:service/:resource_type/actions/:action' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], params['action'])
    return perform_request(params['service'], action, @body)
  end

  # custom resource actions
  post '/:service/:resource_type/:id/actions/:action' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], params['action'])
    @body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, @body)
  end

  # create
  post '/:service/:resource_type' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "create")
    code, hdrs, body = perform_request(params['service'], action, @body)
    #$logger.info "Got headers: #{hdrs}"
    if code == 200 && hdrs['content-type'] == "application/json"
      res = parse_body(body, hdrs['content-type'], "response")
      if res.is_a?(Hash) && res.key?(params['resource_type'].singularize)
        res = res[params['resource_type'].singularize]
      end
      location = "/#{params['service']}/#{params['resource_type']}/#{extract_id(resource, res)}"
      $logger.info "Returning location: #{location}"
      return 201, { "Location" => location }, nil
    elsif code == 201 && hdrs.key?('Location')
      # probably need to massage location header, does this even occur in reality?
      halt 500, "Unimplemented"
      return code, hdrs, body
    elsif code >= 400
      return code, hdrs, body
    elsif code >= 300
      halt 500, "Got a redirect, not yet implemented"
    else
      halt 500, "Got a #{code}, not yet implemented"
    end
  end

  # update
  put '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "update")
    @body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, @body)
  end

  # patch
  patch '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "patch")
    @body.merge!(get_id(resource, params['id']))
    return perform_request(params['service'], action, @body)
  end

  # delete
  delete '/:service/:resource_type/:id' do
    log_info
    svc = get_service(params['service'])
    resource = get_resource(svc, params['service'], params['resource_type'])
    action = get_action(resource, params['resource_type'], "delete")
    @body.merge!(get_id(resource, params['id']))
     code, hdrs, body = perform_request(params['service'], action, @body)
     if code >= 200 && code < 300
       return 204, {}, nil
     else
       return code, hdrs, body
     end
  end

end
