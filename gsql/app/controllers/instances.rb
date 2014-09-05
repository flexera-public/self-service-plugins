  class Instances
    include Praxis::Controller

    implements ApiResources::Instances

    HELLO_WORLD = [ 'Hello world!', 'Привет мир!', 'Hola mundo!', '你好世界!', 'こんにちは世界！' ]

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
