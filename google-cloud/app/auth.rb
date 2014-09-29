class Auth < App

    get '/' do
      if @client
        "Authentication successful"
      else
        services = nil
        services = params[:services].split(',') if params.key?(:services)
        headers 'Content-Type' => 'text/html'
        redirect_url = GoogleCloud.auth_redirect(services)
        redirect redirect_url, <<"EOM"
          <html><body>
            <p>Please visit <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
      end
    end

    get '/redirect' do
      if params[:code] && creds=GoogleCloud.get_creds(params[:code])
        headers 'Content-Type' => 'text/plain'
        "Authentication successful\nCookie=#{creds}"
      else
        #self.response = Praxis::Responses::TemporaryRedirect.new
        #redirect_url = GoogleCloud.auth_redirect(acct, project)
        headers 'Content-Type' => 'text/html'
        #response.headers['Location'] = redirect_url
        <<"EOM"
          <html><body>
            <p>Authentication failed, please retry at
              <a href="#{redirect_url}">Google</a> to authorize</p>
          </body></html>
EOM
      end
    end

end
