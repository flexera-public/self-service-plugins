  class Bds < App
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

__END__

curl -g -XPOST -HContent-Length:0 'http://localhost:9292/tenant/rs-test/bd?bd[name]=rs-test-br2&bd[links][Ctx]=rs-test-net'
curl -g -XDELETE 'http://localhost:9292/tenant/rs-test/bd/rs-test-br2'

