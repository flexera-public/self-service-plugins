class Auth < App

=begin
    get '/'
    def show(acct:, project:, **other_params)
      if GoogleCloudSQL.auth_test(acct)
        response.headers['Content-Type'] = 'text/plain'
        response.body = "Authentication successful"
      else
        self.response = Praxis::Responses::TemporaryRedirect.new
        redirect_url = GoogleCloudSQL.auth_redirect(acct, project)
        response.headers['Content-Type'] = 'text/html'
        response.headers['Location'] = redirect_url
        response.body = <<"EOM"
          <html><body>
            <p>Please visit <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
      end
    end

    post '/:services' do

    def update(acct:, project:, **other_params)
      if GoogleCloudSQL.auth_save(acct, project, request.params.code)
        response.headers['Content-Type'] = 'text/plain'
        response.body = "Your account has been linked to Google Cloud SQL\n"
      else
        self.response = Praxis::Responses::TemporaryRedirect.new
        redirect_url = GoogleCloudSQL.auth_redirect(acct, project)
        response.headers['Content-Type'] = 'text/html'
        response.headers['Location'] = redirect_url
        response.body = <<"EOM"
          <html><body>
            <p>Authentication failed, please retry at
              <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
      end
    end
=end

end
