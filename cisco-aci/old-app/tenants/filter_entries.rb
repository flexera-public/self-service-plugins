  class FilterEntries < App
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
      q.class_filter = "vzEntry"
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
      params[:filter] = params[:Filter] if params[:Filter] && !params[:filter]
      halt 400, "parameter filter_entry missing or not a hash" unless params[:filter_entry] && params[:filter_entry].is_a?(Hash)
      halt 400, "filter_entry hash must have a name field" unless params[:filter_entry]['name']
      puts "FilterEntry is: #{params[:filter_entry].inspect}"
      halt 400, "parameter filter missing" unless params[:filter] && params[:filter].is_a?(String)
      params[:filter_entry].delete("deployment_href")
      params[:filter_entry].delete("filter")
      params[:filter_entry]["etherT"] = "ip"
      filter_name = params[:filter].sub(/^.*\//, '')

      uni = ACIrb::PolUni.new(nil)
      tenant = ACIrb::FvTenant.new(uni, name: params[:tenant])
      filter = ACIrb::VzFilter.new(tenant, name: filter_name)
      filter_entry = ACIrb::VzEntry.new(filter)
      name = params[:filter_entry]['name']
      add_stuff(filter_entry, params[:filter_entry])
      puts "Filter obj is: #{filter_entry.inspect}"
      begin
        result = filter_entry.create(@api)
        puts "ACI returned #{result.inspect}"
        [ 201, { 'Location' => "/tenants/#{params[:tenant]}/filters/#{filter_name}/filter_entrys/#{name}" }, "" ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

  end

__END__

curl -g -XPOST -HContent-Length:0 'http://localhost:9292/tenants/rs-test/filter_entry?filter_entry[name]=http-entry&filter=web-filter&filter_entry[dFromPort]=80&filter_entry[dToPort]=80&filter_entry[prot]=tcp&filter_entry[etherT]=ip'
curl -g -XDELETE 'http://localhost:9292/tenants/rs-test/filter_entry/http-entry'

