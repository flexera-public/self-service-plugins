  class Instances
    include Praxis::Controller
    implements ApiResources::Instances

    include GoogleCloudSQLMixin

    SAMPLE = [
      { id: "i-1234", state: "running", region: "us-central" },
      { id: "i-1235", state: "inactive", region: "us-central" },
    ]

    def index(acct:, **params)
      result = @gc_sql_client.execute(
        api_method: @gc_sql_api.instances.list,
        parameters: { project: @gc_sql_project },
      )
      puts "Got #{result.status.inspect}"
      if result.success?
        Praxis::Responses::Ok.new(
          headers: { 'Content-Type' => 'vnd.rightscale.instance+json;type=collection' },
          body: SAMPLE
        )
      else
        puts "Error: #{result.inspect}"
        Praxis::Responses::BadRequest.new(
          headers: { 'Content-Type' => 'vnd.rightscale.instance+json;type=collection' },
          body: "#{result.error_message} The request was: #{result.request.inspect}",
        )
      end
    end

    def show(acct:, id:, **other_params)
      if id.to_i < 2
        response.body = SAMPLE[id.to_i]
      else
        response.status = 404
        response.body   = { error: '404: Not found' }
      end
      response.headers['Content-Type'] = 'vnd.rightscale.instance+json;type=collection'
      response
    end

  end
