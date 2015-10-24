  class Bd < App
    require 'acirb'

    before do
      @api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
                                   format: "json", debug: false)
      @api.debug = true
    end

    helpers do

      def add_stuff(obj, stuff)
        links = stuff.delete('links')
        stuff.each_pair do |k,v|
          if obj.props.key?(k)
            obj.set_prop(k, v)
          else
            halt 400, "Oops: #{obj.class_name} does not have attribute #{k}, valid attributes: #{obj.props.keys.sort.join(' ')}"
          end
        end
        if links.is_a?(Hash)
          links.each do |k, v|
            puts "Link: #{k}->#{v}"
            child = nil
            begin
              child = Object.const_get("ACIrb::Fv#{k}").new(obj)
              #obj.add_child(child)
            rescue
              begin
                child = Object.const_get("ACIrb::FvRs#{k}").new(obj)
                #obj.add_child(child)
              rescue
                halt 400, "Cannot create link '#{k}' for '#{obj.class}'"
              end
            end
            if child.props.key?("name")
              child.set_prop('name', v)
            else
              name_props = child.props.select{|p,v| p.end_with?('Name')}
              if name_props.size == 1
                puts "Setting prop #{name_props.first[0]}=#{v} (#{name_props.inspect})"
                child.set_prop(name_props.first[0], v)
              else
                halt 400, "Cannot set name for link '#{k}': #{name_props.sort.join(' ')}"
              end
            end
            puts "Child #{k}:#{v} is: #{child.inspect}"
          end
          stuff.delete('links')
        end
        obj
      end

      def gen_json(obj)
        if obj.is_a?(Array)
          obj.map{|o|o.to_json}
        else
          obj.to_json
        end
      end

    end

    get '/' do
      q = ACIrb::DnQuery.new("uni/tn-#{params[:tenant]}")
      q.class_filter = "fvBD"
      q.query_target = 'children'
      result = @api.query(q)
      puts "ACI returned #{result.inspect}"
      if result
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      else
        puts "Error: #{result.inspect}"
        halt result.status, result.error_message
      end
    end

    post '/' do
      halt 400, "parameter bd missing or not a hash" unless params[:bd] && params[:bd].is_a?(Hash)
      halt 400, "bd hash must have a name field" unless params[:bd]['name']
      puts "BD is: #{params[:bd].inspect}"

      uni = ACIrb::PolUni.new(nil)
      tenant = ACIrb::FvTenant.new(uni, name: params[:tenant])
      bd = ACIrb::FvBD.new(tenant)
      add_stuff(bd, params[:bd])
      puts "BD obj is: #{bd.inspect}"
      begin
        result = bd.create(@api)
        puts "ACI returned #{result.inspect}"
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

    delete '/:bd' do
      uni = ACIrb::PolUni.new(nil)
      tenant = ACIrb::FvTenant.new(uni, name: params[:tenant])
      bd = ACIrb::FvBD.new(tenant, name: params[:bd])
      begin
        result = bd.destroy(@api)
        puts "ACI returned #{result.inspect}"
        [ 204, { 'Content-Type' => 'application/json' }, "" ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

  end
