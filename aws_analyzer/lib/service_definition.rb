module Analyzer

  # Generic service definition that can be used to drive plugin services
  class ServiceDefinition

    # [String] Service name
    attr_reader :name

    # [String] Service url
    attr_reader :url

    # [String] API Version
    attr_reader :version

    # [Hash] Cloud provider specific metadata needed by API client
    attr_reader :metadata

    # [Array<Resource>] List of resources exposed by service
    attr_reader :resources

    # [Array<Shapes>] Request and response body structure descriptions
    attr_reader :shapes

    # Setup using given hash
    def initialize(opts)
      opts.each { |k, v| instance_variable_set("@#{k}", v) }
    end

    # Convert to hash
    def to_hash
      { name:      @name,
        url:       @url,
        version:   @version || '1.0',
        metadata:  @metadata,
        resources: @resources.map(&:to_hash),
        shapes:    @shapes }
    end

    # Is definition empty?
    def empty?
      @resources.nil? || @resources.empty?
    end

    # YAML conversion
    def to_yaml
      YAML.dump(to_hash)
    end
    alias :to_s :to_yaml

  end

end
