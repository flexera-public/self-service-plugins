require_relative 'service'
require_relative 'resource'
require_relative 'media_type'

module Contexts
  module DSL
    def create_service(name)
      const_set(
        name.camel_case,
        Contexts::Service.create
      )
    end

    def create_resource(service, name)
      service.const_set(
        name.camel_case,
        Contexts::Resource.create
      )
    end

    def create_media_type(service, name)
      service.const_set(
        name.camel_case,
        Contexts::MediaType.create
      )
    end
  end
end
