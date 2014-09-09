  class Auth
    include Praxis::Controller
    implements ApiResources::Auth

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

    def show(acct:, **other_params)
      if GoogleCloudSQL.auth_test(acct)
        response.headers['Content-Type'] = 'text/plain'
        response.body = "Authentication successful"
      else
        response = Praxis::Responses::TemporaryRedirect.new
        redirect_url = GoogleCloudSQL.auth_redirect(acct)
        response.headers['Content-Type'] = 'text/html'
        response.headers['Location'] = redirect_url
        response.body = <<"EOM"
          <html><body>
            <p>Please visit <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
      end
      response
    end

    def update(acct:, **other_params)
      if request.params.code
        if GoogleCloudSQL.auth_set(acct, request.params.code)
          response.headers['Content-Type'] = 'text/plain'
          response.body = "Your account has been linked to Google Cloud SQL\n"
        else
          response = Praxis::Responses::TemporaryRedirect.new
          redirect_url = GoogleCloudSQL.auth_redirect(acct)
          response.headers['Content-Type'] = 'text/html'
          response.headers['Location'] = redirect_url
          response.body = <<"EOM"
            <html><body>
              <p>Authentication failed, please retry at
                <a href="#{redirect_url}">Google</a> to authorize</p>
            </body></html>
EOM
        end
      else
        response = Praxis::Responses::BadRequest.new
        response.headers['Content-Type'] = 'text/plain'
        response.body = "Authentication code missing\n"
      end
      response
    end

  end
