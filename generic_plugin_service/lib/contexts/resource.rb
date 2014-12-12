module Contexts
  class Resource
    def self.create
      Class.new do
        include Praxis::ResourceDefinition
      end
    end
  end
end
