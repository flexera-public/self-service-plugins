  class Mo < App
    require 'acirb'

    #before do
    #  @api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
    #                               format: "json", debug: false)
    #  @api.debug = true
    #end

    helpers do

      # construct an resource object of the correct class with all the parent resources
      # instantiated as well, returns resource object and it's canonical href
      def make_resource_path(context)
        # start with the policy universe root and traverse down to the resource requested
        node = ACIrb::PolUni.new(nil)
        href = "/mo"
        context.each_slice(2) do |kl, name|
          kl = kl.sub(/ies$/,'y').sub(/s$/,'') # singularize hack
          #puts "Adding: #{kl} / #{name}"
          # find a child class of the current node that matches what's asked for
          # first try exact match, then try suffix match (allows user to omit leading Fv/Vz/...
          child_classes = node.child_classes.sort
          klass_name = child_classes.select{|k| k.to_s == kl}
          klass_name = child_classes.select{|k| k.to_s.end_with?(kl)} if klass_name.size == 0
          klass_name.sort!{|a,b| a.length <=> b.length} if klass_name.size > 1
          #halt 400, "Ambiguous class #{kl}. Matching child classes of #{node.ruby_class} are: "+
          #  klass_name.join(' ') if klass_name.size>1
          halt 404, "Unknown class #{kl}. Child classes of #{node.ruby_class} are: "+
            child_classes.join(' ') if klass_name.size == 0
          klass_name = klass_name.first
          #puts "Selected #{klass_name}"
          klass = Object.const_get("ACIrb::#{klass_name}")
          # allocate new object with parent set
          node = klass.new(node, {klass.naming_props.first.to_sym => name})
          # also add to href
          href += "/#{klass_name}/#{name}"
        end
        [node, href]
      end

    end

    # INDEX resources with filter
    get %r{\A((/[\w]+/[^/]+)*/[\w]+)\z} do |*ctx| # match (/class/:id)*/class
      puts "***** INDEX ***** INDEX *****"

      halt 400, "parameter parent must be a string" \
        if params[:parent] && !params[:parent].is_a?(String)

      halt 400, "parameter filter must be an array" \
        unless params[:filter] && params[:filter].is_a?(Array)
      filters = params[:filter].map{|f| f.split('==')}.to_h
      halt 400, "filter must include 'name==' expression" unless filters.key?("name")

      ctx = ctx[0]
      if filters['parent']
        filters['parent'].sub!(/^\w+: /, '') # garbage added by v1 plugin interface
        ctx = filters['parent'].sub(%r{^/mo},'') + ctx.sub(%r{^/(Fv)?Tenant/[^/]+},'')
      end
      ctx += '/' + filters['name']
      ctx = ctx.sub(%r{^/}, "").split('/')
      puts "Route match: #{ctx}"
      node, href = make_resource_path(ctx)

      # fetch and return
      begin
        result = $api.get(url: "/api/mo/#{node.dn}.#{$api.format}")
        $logger.debug "ACI returned #{result.inspect}"
        halt 404, "#{href} not found" if result.size == 0
        result = result[0].attributes
        result["href"] = href
        $logger.info "Returning #{result.inspect}"
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end

    end

    # GET a resource
    get %r{\A((/[\w]+/[^/]+)+)\z} do |*ctx| # match (/class/:id)+
      puts "Route match: #{ctx[0].inspect}"
      ctx = ctx[0].sub(%r{^/}, "").split('/')
      node, href = make_resource_path(ctx)

      # fetch and return
      begin
        result = $api.get(url: "/api/mo/#{node.dn}.#{$api.format}")
        #puts "ACI returned #{result.inspect}"
        halt 404, "#{href} not found" if result.size == 0
        result = result[0].attributes
        result["href"] = href
        $logger.info "Returning #{result.inspect}"
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

    post %r{\A((/[\w]+/[^/]+)*/[\w]+)\z} do |*ctx| # match (/class/:id)*/class
      halt 400, "parameter props must be a hash" \
        if params[:props] && !params[:props].is_a?(Hash)
      halt 400, "parameter parent must be a string" \
        if params[:parent] && !params[:parent].is_a?(String)
      if params[:props]
        # some clean-up due to cloud plugins V1 shtuff
        params[:props].delete("deployment_href")
        params[:parent] = params[:props]["parent"] if params[:props]["parent"]
        params[:props].delete("parent")
      end

      ctx = ctx[0]
      if params[:parent]
        ctx = params[:parent].sub(%r{^/mo},'') + ctx.sub(%r{^/(Fv)?Tenant/[^/]+},'')
      end
      ctx += '/' # append empty name property
      ctx = ctx.sub(%r{^/}, "").split('/')
      puts "Route match: #{ctx}"
      node, href = make_resource_path(ctx)

      # populate with properties
      add_stuff(node, params[:props])
      # Extract the name for the href
      name = node.attributes[node.naming_props.first]
      href += name

      puts "Resource is: #{node.inspect}"
      puts "Resource href: #{href}"
      begin
        result = node.create($api)
        puts "ACI returned #{result.inspect}"
        [ 201, { 'Location' => href }, "" ]
      rescue ACIrb::RestClient::ApicErrorResponse => e
        puts "Error: #{e.message}"
        halt 500, e.message
      end
    end

    delete %r{\A((/[\w]+/[^/]+)+)\z} do |*ctx| # match (/class/:id)+
      puts "Route match: #{ctx[0].inspect}"
      ctx = ctx[0].sub(%r{^/}, "").split('/')
      node, href = make_resource_path(ctx)

      # delete and return
      begin
        result = node.destroy($api)
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

