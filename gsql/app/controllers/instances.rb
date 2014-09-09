  class Instances
    include Praxis::Controller
    implements ApiResources::Instances
    include GoogleCloudSQLMixin

=begin
    before :action do |controller|
#     puts methods.sort.join(" ")
      puts "Auth befpore filter"
      acct = controller.request.params.acct

      raise "Authentication is missing" unless acct
      @gc_sql_client = GoogleCloudSQL.client(acct)
      raise "Authentication failed" unless @gc_sql_client
      @gc_sql_api = GoogleCloudSQL.api
      raise "Internal error: cannot retrieve Cloud SQL API definition" unless @gc_sql_api
    end
=end

    before :action do |controller|
      puts "Instances before action"
    end

    SAMPLE = [
      { id: "i-1234", state: "running", region: "us-central" },
      { id: "i-1235", state: "inactive", region: "us-central" },
    ]

    def index(acct:, **params)
      response.headers['Content-Type'] = 'vnd.rightscale.instance+json;type=collection'
      response.body = SAMPLE
      response
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
