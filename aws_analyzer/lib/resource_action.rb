module Analyzer

  # A service resource action
  class ResourceAction

    # [String] Action name (e.g. "create")
    attr_reader :name

    # [Symbol] Action verb (one of :get, :post, :put, :delete, ...)
    attr_reader :verb

    # [String] Action path
    attr_reader :path

    # [Shape] Action payload
    attr_reader :payload

    # [Array<Shape>] Action params
    attr_reader :params

    # [Shape] Action response
    attr_reader :response

    # Initialize
    def initialize
    end

    # Hash representation
    def to_hash
      { name:     @name,
        verb:     @verb,
        path:     @path,
        payload:  @payload.name,
        params:   @params.map(&:name),
        response: @response.name }
    end

    # YAML representation
    def to_yaml
      YAML.dump(to_hash)
    end
    alias :to_s :to_yaml

  end

end
