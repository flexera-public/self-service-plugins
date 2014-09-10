  class Buckets
    include Praxis::Controller
    implements ApiResources::Buckets

    include GoogleCloudStorageMixin

    def make_href(acct, id)
      #puts "Href for acct=#{acct} bucket=#{id}"
      route = ApiResources::Buckets.actions[:show].named_routes[:bucket_href]
      route.path.expand(acct: acct, id: id)
    end

    # Convert the google cloud storage bucket representation to "our" representation
    def convert_bucket(acct, b)
      b['href'] = make_href(acct, b['id'])
      b
    end

    def index(acct:, **params)
      result = @gc_storage_client.execute(
        api_method: @gc_storage_api.buckets.list,
        parameters: { project: @gc_storage_project },
      )
      puts "Google returned #{result.status.inspect}"
      puts "Google returned #{result.data.inspect}"
      if result.success? && result.data?
        Praxis::Responses::Ok.new(
          headers: { 'Content-Type' => 'vnd.rightscale.bucket+json;type=collection' },
          body: MultiJson.load(result.body)['items'].
                  select{|i| i['kind'] == "storage#bucket"}.
                  map{|i| convert_bucket(acct, i)},
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
      result = @gc_storage_client.execute(
        api_method: @gc_storage_api.buckets.get,
        parameters: { project: @gc_storage_project, bucket: id },
      )
      puts "Google returned #{result.status.inspect}"
      if result.success? && result.data?
        Praxis::Responses::Ok.new(
          headers: { 'Content-Type' => 'vnd.rightscale.bucket+json;type=item' },
          body: convert_bucket(acct, MultiJson.load(result.body)),
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
      b = request.raw_params['b']
      result = @gc_storage_client.execute(
        api_method: @gc_storage_api.buckets.insert,
        parameters: { project: @gc_storage_project },
        body_object: b,
      )
      puts "Google returned #{result.status.inspect}"
      if result.success?
        Praxis::Responses::Created.new(
          headers: { 'Location' => make_href(acct, b['name']) },
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
      result = @gc_storage_client.execute(
        api_method: @gc_storage_api.buckets.delete,
        parameters: { project: @gc_storage_project, bucket: id },
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
