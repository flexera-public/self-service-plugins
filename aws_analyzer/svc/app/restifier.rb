  class Restifier < App
    $path = "/home/src/aws-sdk-core-ruby/aws-sdk-core/apis"

    $services = {}

    helpers do
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

    end


  get '/:service/:resource_type' do
    $logger.info "Params   : #{params.inspect}"
    body_str = request.body.is_a?(String) ? request.body : request.body.read
    #$logger.info "Body     : #{body_str}"
    $logger.info "Body     : #{request.content_type} (#{body_str.size} bytes)"
    $logger.info "Service  : #{params['service']}"
    $logger.info "Res Type : #{params['resource_type']}"

    # Get the metadata for the service and make sure the operation exists
    svc = get_service(params['service'])
    $logger.info "Metadata: #{svc.metadata.inspect}"

    unless client.operation_names.include?(params['operation'].to_sym)
      halt 404, "Operation #{params['operation']} is not supported by service #{params['service']}"
    end

    # Get the body of the request and make it ready to be the operation's parameter
    body = {}
    if body_str.size > 0
      unless request.content_type && request.content_type.start_with?("application/json")
        halt 400, "Request must contain application/json parameters"
      end
      begin
        body = Yajl::Parser.parse(request.body, :symbolize_keys => true)
      rescue StandardError => e
        halt 400, "Error parsing json body: #{e}"
      end
      unless body.is_a?(Hash)
        halt 400, "Request body must consist of a json hash"
      end
      $logger.info "Request body contains: #{body.keys.sort.join(' ')}"
    end

    # Perform the operation
    begin
      res = client.send(params['operation'], body)
    rescue Aws::Errors::ServiceError => e
      code = e.context.http_response.status_code
      #$logger.info "*** Service error: #{e.context.http_response.inspect}"
      $logger.info "*** AWS returned error: #{code} \"#{e}\""
      halt code, e.message
    rescue Aws::Errors, ArgumentError => e
      $logger.info "*** AWS gem error: #{e}"
      halt 400, e.message
    rescue Exception => e
      $logger.info "*** Error: #{e} #{e.inspect}"
      $logger.info e.backtrace[0..1].join(' | ')
      halt 400, e.message
    end

    if res.is_a?(Aws::PageableResponse)
      $logger.info "Pageable response"

      if res.last_page?
        code = res.context.http_response.status_code
        #$logger.info "Result: #{res.data.inspect}"
        #$logger.info "Result code: #{res.context.http_response.status_code}"
        if res.data.is_a?(Struct)
          return code, { 'content-type' => "application/json" }, Yajl::Encoder.encode(res.data.to_h)
        else
          return code, { 'content-type' => "application/json" }, Yajl::Encoder.encode(res.data)
        end
      end

      response_key = nil
      response_array = []
      code = 400
      res.each_page do |page|
        data = page.data.to_h
        #$logger.info "Data: #{data.inspect}"
        #$logger.info "Page: #{page.inspect}"
        if !response_key
          # find something in the response that is an array
          data.each_pair do |k, v|
            if v.is_a?(Array)
              response_key = k
              break
            end
          end
          #data.each_pair do |k, v|
          #  $logger.info "Data[#{k}] -> #{v.class}"
          #end
          if !response_key
          halt 500, "Cannot locate an array in AWS response having #{data.keys.join(' ')}"
          end
        end
        # now add what we got to the total
        response_array += data[response_key]
        code = page.context.http_response.status_code
      end
      response = Yajl::Encoder.encode(response_key => response_array)
      $logger.info "Returning #{response_array.size} elements as #{response.size} bytes"

      return code, { 'content-type' => "application/json" }, response
    else
      $logger.info "Result: #{res.inspect}"
      $logger.info "Result: #{res.data.inspect}"
      halt 500, "result is not pageable, dunno what to do"
    end
  end

  end
