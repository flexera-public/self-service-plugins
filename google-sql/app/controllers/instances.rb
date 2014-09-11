  class Instances
    include Praxis::Controller
    implements ApiResources::Instances

    include GoogleCloudSQLMixin

    def make_href(acct, id)
      route = ApiResources::Instances.actions[:show].named_routes[:instance_href]
      route.path.expand(acct: acct, id: id)
    end

    # Convert the google cloud sql instance representation to "our" representation
    def convert_instance(acct, i)
      route = ApiResources::Instances.actions[:show].named_routes[:instance_href]
      i['href'] = make_href(acct, i['instance'])
      i
    end

    def index(acct:, **params)
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.list,
        parameters: { project: @gc_sql_project },
      )
      puts "Google returned #{result.status.inspect}"
      if result.success? && result.data?
        Praxis::Responses::Ok.new(
          headers: { 'Content-Type' => 'vnd.rightscale.instance+json;type=collection' },
          body: MultiJson.load(result.body)['items'].
                  select{|i| i['kind'] == "sql#instance"}.
                  map{|i| convert_instance(acct, i)},
        )
      else
        puts "Error: #{result.inspect}"
        Praxis::Responses::BadRequest.new(
          headers: { 'Content-Type' => 'text/plain' },
          body: "#{result.error_message}", # The request was: #{result.request.inspect}",
        )
      end
    end

    def show(acct:, id:, **other_params)
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.get,
        parameters: { project: @gc_sql_project, instance: id },
      )
      puts "Google returned #{result.status.inspect}"
      if result.success? && result.data?
        Praxis::Responses::Ok.new(
          headers: { 'Content-Type' => 'vnd.rightscale.instance+json;type=item' },
          body: convert_instance(acct, MultiJson.load(result.body)),
        )
      else
        puts "Error: #{result.inspect}"
        Praxis::Responses::BadRequest.new(
          headers: { 'Content-Type' => 'text/plain' },
          body: "#{result.error_message}", # The request was: #{result.request.inspect}",
        )
      end
    end

    def create(acct:, **other_params)
      #i = JSON.parse(request.raw_payload)
      #i = request.raw_params['i']
      #i = request.raw_params
      i = JSON.parse(request.raw_payload)
      i['settings'] ||= {}
      i['settings']['tier'] = i['tier']
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.insert,
        parameters: { project: @gc_sql_project },
        body_object: i,
      )
      puts "Google returned #{result.status.inspect}"
      if result.success?
        Praxis::Responses::Created.new(
          headers: { 'Location' => make_href(acct, i['instance']) },
        )
      else
        puts "Error: #{result.inspect}"
        Praxis::Responses::BadRequest.new(
          headers: { 'Content-Type' => 'text/plain' },
          body: "#{result.error_message}", # The request was: #{result.request.inspect}",
        )
      end
    end

    def delete(acct:, id:, **other_params)
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.delete,
        parameters: { project: @gc_sql_project, instance: id },
      )
      puts "Google returned #{result.status.inspect}"
      if result.success?
        Praxis::Responses::NoContent.new
      else
        puts "Error: #{result.inspect}"
        Praxis::Responses::BadRequest.new(
          headers: { 'Content-Type' => 'text/plain' },
          body: "#{result.error_message}", # The request was: #{result.request.inspect}",
        )
      end
    end

  end
