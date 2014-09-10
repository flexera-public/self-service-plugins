module FilterTest
  def self.included(klass)
    klass.class_eval do
      before :action do |controller|
        puts "*** FilterTest before filter"
      end
    end
  end
end

module V1
  class Hello
    include Praxis::Controller
    implements V1::ApiResources::Hello
    include FilterTest

    HELLO_WORLD = [ 'Hello world!', 'Привет мир!', 'Hola mundo!', '你好世界!', 'こんにちは世界！' ]

    before :action do |controller|
      puts "*** Internal before action"
      #puts "Self is a #{self.inspect}"
      #puts "Self.methods #{self.methods.sort.join(' ')}"
      #puts "Self.private_methods #{self.private_methods.sort.join(' ')}"
      #puts "Self.class_variables #{self.class_variables.sort.join(' ')}"
      #puts "Self.instance_variables #{self.instance_variables.sort.join(' ')}"
      #puts "The request is #{controller.request.inspect}"
    end

    def index(**params)
      response.headers['Content-Type'] = 'application/json'
      response.body = HELLO_WORLD.to_json
      response
    end

    def show(id:, **other_params)
      hello = HELLO_WORLD[id]
      if hello
        response.body = { id: id, data: hello }
      else
        response.status = 404
        response.body   = { error: '404: Not found' }
      end
      response.headers['Content-Type'] = 'application/json'
      response
    end
  end
end
