  class Auth
    include Praxis::Controller
    implements ApiResources::Auth

    def show(acct:, project:, **other_params)
      if GoogleCloudSQL.auth_test(acct)
        resp = Praxis::Responses::Ok.new
        resp.headers['Content-Type'] = 'text/plain'
        resp.body = "Authentication successful"
        resp
      else
        resp = Praxis::Responses::TemporaryRedirect.new
        redirect_url = GoogleCloudSQL.auth_redirect(acct, project)
        resp.headers['Content-Type'] = 'text/html'
        resp.headers['Location'] = redirect_url
        resp.body = <<"EOM"
          <html><body>
            <p>Please visit <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
        resp
      end
    end

    def update(acct:, project:, **other_params)
      if GoogleCloudSQL.auth_save(acct, project, request.params.code)
        resp = Praxis::Responses::Ok.new
        resp.headers['Content-Type'] = 'text/plain'
        resp.body = "Your account has been linked to Google Cloud SQL\n"
        resp
      else
        resp = Praxis::Responses::TemporaryRedirect.new
        redirect_url = GoogleCloudSQL.auth_redirect(acct, project)
        resp.headers['Content-Type'] = 'text/html'
        resp.headers['Location'] = redirect_url
        resp.body = <<"EOM"
          <html><body>
            <p>Authentication failed, please retry at
              <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
        resp
      end
    end

  end
