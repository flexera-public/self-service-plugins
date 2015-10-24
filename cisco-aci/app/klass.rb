  class Klass < App
    require 'acirb'

    before do
      @api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
                                   format: "json", debug: true)
    end

    helpers do

      def gen_json(obj)
        if obj.is_a?(Array)
          obj.map{|o|o.to_json}
        else
          obj.to_json
        end
      end

=begin
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
=end

    end

=begin
    get '/' do
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.list,
        parameters: { project: @gc_sql_project },
      )
      puts "Google returned #{result.status.inspect}"

      if result.success? && result.data?
        [ 200, { 'Content-Type' => 'application/json' },
          MultiJson.load(result.body)['items'].
            select{|i| i['kind'] == "sql#instance"}.
            map{|i| convert_instance(acct, i)} ]
      else
        puts "Error: #{result.inspect}"
        halt result.status, result.error_message
      end
    end
=end

    get '/:klass' do
      result = @api.lookupByClass(params[:klass])
      puts "ACI returned #{result.inspect}"
      [ 200, { 'Content-Type' => 'application/json' }, gen_json(result)]

=begin
      if result.success? && result.data?
        [ 200, { 'Content-Type' => 'application/json' },
          convert_instance(acct, MultiJson.load(result.body)) ]
      else
        puts "Error: #{result.inspect}"
        halt result.status, result.error_message
      end
=end
    end

=begin
    post '/' do
      #i = JSON.parse(request.raw_payload)
      #i = request.raw_params['i']
      #i = request.raw_params
      i = params[:i]
      i['settings'] ||= {}
      i['settings']['tier'] = i.delete('tier')
      puts "Payload: #{i.inspect}"
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.insert,
        parameters: { project: @gc_sql_project },
        body_object: i,
      )
      puts "Google returned #{result.status.inspect}"

      if result.success?
        [ 201, { 'Location' => make_href(acct, i['instance']) }, nil ]
      else
        puts "Error: #{result.inspect}"
        halt result.status, result.error_message
      end
    end

    delete ':id' do
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.delete,
        parameters: { project: @gc_sql_project, instance: params[:id] },
      )
      puts "Google returned #{result.status.inspect}"

      if result.success?
        [ 204, {}, nil ]
      else
        puts "Error: #{result.inspect}"
        halt result.status, result.error_message
      end
    end
=end

  end
