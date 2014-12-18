module Analyzer

  # A service resource action
  # This data structure is common to all clouds
  class Action

    # [String] Action name (e.g. "create_stack")
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

    # Initialize with hash
    def initialize(opts)
      opts.each { |k, v| instance_variable_set("@#{k}", v) }
    end

    # Hash representation
    def to_hash
      { 'name'     => @name,
        'verb'     => @verb,
        'path'     => @path,
        'payload'  => @payload,
        'params'   => @params,
        'response' => @response }
    end

    # YAML representation
    def to_yaml
      YAML.dump(to_hash)
    end
    alias :to_s :to_yaml

  end

end
