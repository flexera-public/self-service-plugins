  class Mo < App
    require 'acirb'

    helpers do

      # construct a resource object of the correct class with all the parent resources
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

      # add parameters as properties or child relationship to an object
      def add_stuff(obj, stuff)
        #$logger.debug "Add stuff #{obj.class_name}: #{stuff.inspect}"
        stuff.each_pair do |k,v|
          # if it's a property, just set it
          if obj.props.key?(k)
            obj.set_prop(k, v)
            next
          end

          # see whether it's a child relationship resource
          cap_k = k[0].capitalize + k[1,k.length-1]
          cc = obj.child_classes.select{|cc| cc == cap_k || cc =~ /Rs#{cap_k}\z/} # also Rt?
          halt 400, "Ambiguous child class #{k} in #{obj.ruby_class}, choices: #{cc.sort.join(' ')}" \
            if cc.size > 1
          v.sub!(%r{^/mo.*/}, '') # convert value from href to name
          if cc.size == 1
            cc = cc.first
            child = Object.const_get("ACIrb::#{cc}").new(obj)
            if child.props.key?("name")
              child.set_prop('name', v)
            else
              name_props = child.props.select{|p,v| p.end_with?('Name')}
              if name_props.size == 1
                #$logger.debug "Setting prop #{name_props.first[0]}=#{v} (#{name_props.inspect})"
                child.set_prop(name_props.first[0], v)
              else
                halt 400, "Cannot set name for link '#{k}': #{name_props.sort.join(' ')}"
              end
            end
            next
          end

          halt 400, "Oops: #{obj.class_name} does not have attribute or child class #{k},\n" +
            "valid attributes: #{obj.props.keys.sort.join(' ')},\n" +
            "valid child classes: #{obj.child_classes.sort.join(' ')}"
        end
        obj
      end

      # convert an object, a hash, or an array to JSON
      def gen_json(obj)
        if obj.is_a?(Array)
          obj.map{|o|o.to_json}
        else
          obj.to_json
        end
      end

      # run a block that accesses the ACI API and rescue exceptions, retry if the auth token
      # has expired
      def aci_op
        count = 0
        begin
          yield
        rescue ACIrb::RestClient::ApicErrorResponse => e
          puts e.message
          if count < 1 && e.message =~ /Error: Token timeout/
            $api = ACIrb::RestClient.new(url: $apic_url, user: $username, password: $password,
                                         format: "json", debug: false)
            count += 1
            puts "Retrying"
            retry
          end
          #halt 400, e.message
          [ 400, e.message ]
        end
      end

    end

    # INDEX resources with filter
    get %r{[\A]((/[\w]+/[^/]+)*/[\w]+)} do |*ctx| # match (/class/:id)*/class
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
      aci_op do
        result = $api.get(url: "/api/mo/#{node.dn}.#{$api.format}")
        $logger.debug "ACI returned #{result.inspect}"
        halt 404, "#{href} not found" if result.size == 0
        result = result[0].attributes
        result["href"] = href
        $logger.info "Returning #{result.inspect}"
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      end
    end

    # GET a resource
    get %r{[\A]((/[\w]+/[^/]+)+)} do |*ctx| # match (/class/:id)+
      puts "Route match: #{ctx[0].inspect}"
      ctx = ctx[0].sub(%r{^/}, "").split('/')
      node, href = make_resource_path(ctx)

      # fetch and return
      aci_op do
        result = $api.get(url: "/api/mo/#{node.dn}.#{$api.format}")
        #puts "ACI returned #{result.inspect}"
        halt 404, "#{href} not found" if result.size == 0
        result = result[0].attributes
        result["href"] = href
        $logger.info "Returning #{result.inspect}"
        [ 200, { 'Content-Type' => 'application/json' }, gen_json(result) ]
      end
    end

    # CREATE a resource
    post %r{[\A]((/[\w]+/[^/]+)*/[\w]+)} do |*ctx| # match (/class/:id)*/class
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
      aci_op do
        result = node.create($api)
        puts "ACI returned #{result.inspect}"
        [ 201, { 'Location' => href }, "" ]
      end
    end

    # DELETE a resource
    delete %r{[\A]((/[\w]+/[^/]+)+)} do |*ctx| # match (/class/:id)+
      puts "Route match: #{ctx[0].inspect}"
      ctx = ctx[0].sub(%r{^/}, "").split('/')
      node, href = make_resource_path(ctx)

      # delete and return
      aci_op do
        result = node.destroy($api)
        puts "ACI returned #{result.inspect}"
        [ 204, { 'Content-Type' => 'application/json' }, "" ]
      end
    end

  end

__END__

curl -g -XPOST -HContent-Length:0 'http://localhost:9292/tenants/rs-test/filter?filter[name]=web-filter'
curl -g -XDELETE 'http://localhost:9292/tenants/rs-test/filter/web-filter'
