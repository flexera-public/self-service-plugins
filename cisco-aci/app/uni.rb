  class Uni < App
    require 'acirb'

    before do
      @api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
                                   format: "raw", debug: false)
      @api.debug = true
    end

    helpers do

      def gen_json(obj)
        if obj.is_a?(Array)
          obj.map{|o|o.to_json}
        else
          obj.to_json
        end
      end

    end

    post '/' do
      body = request.body.read
      #puts "Posting:", body
      puts "Debug: #{@api.debug}"
      result = @api.post(url: '/api/mo/uni.xml', data: body)
      puts "ACI returned #{result.inspect}"

      if result
        [ 200, { 'Content-Type' => 'application/json' }, result ]
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

  end
