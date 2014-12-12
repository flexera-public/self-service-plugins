module Contexts
  class MediaType
    def self.create
      Class.new(Praxis::MediaType)
    end
  end
end

