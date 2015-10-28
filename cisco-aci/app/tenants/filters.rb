  class Filters < App
    require 'acirb'

    before do
      @api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
                                   format: "json", debug: false)
      @api.debug = true
    end

    helpers do
    end

    get '/' do
      q = ACIrb::DnQuery.new("uni/tn-#{params[:tenant]}")
      q.class_filter = "vzFilter"
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

    get '/:filter' do
      uni = ACIrb::PolUni.new(nil)
      tenant = ACIrb::FvTenant.new(uni, name: params[:tenant])
      filter = ACIrb::VzFilter.new(tenant, name: params[:filter])
      begin
        result = @api.get(url: "/api/mo/#{filter.dn}.#{@api.format}")
        puts "ACI returned #{result.inspect}"
        result = result[0].attributes
        result["href"] = "/tenants/#{params[:tenant]}/filters/#{result["name"]}"
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

    post '/' do
      params[:filter] = params[:Filter] if params[:Filter] && !params[:filter]
      halt 400, "parameter filter missing or not a hash" unless params[:filter] && params[:filter].is_a?(Hash)
      halt 400, "filter hash must have a name field" unless params[:filter]['name']
      puts "Filter is: #{params[:filter].inspect}"
      params[:filter].delete("deployment_href")

      uni = ACIrb::PolUni.new(nil)
      tenant = ACIrb::FvTenant.new(uni, name: params[:tenant])
      filter = ACIrb::VzFilter.new(tenant)
      name = params[:filter]['name']
      add_stuff(filter, params[:filter])
      puts "Filter obj is: #{filter.inspect}"
      begin
        result = filter.create(@api)
        puts "ACI returned #{result.inspect}"
        [ 201, { 'Location' => "/tenants/#{params[:tenant]}/filters/#{name}" }, "" ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

    delete '/:filter' do
      uni = ACIrb::PolUni.new(nil)
      tenant = ACIrb::FvTenant.new(uni, name: params[:tenant])
      filter = ACIrb::VzFilter.new(tenant, name: params[:filter])
      begin
        result = filter.destroy(@api)
        puts "ACI returned #{result.inspect}"
        [ 204, { 'Content-Type' => 'application/json' }, "" ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

  end

__END__

curl -g -XPOST -HContent-Length:0 'http://localhost:9292/tenants/rs-test/filter?filter[name]=web-filter'
curl -g -XDELETE 'http://localhost:9292/tenants/rs-test/filter/web-filter'

