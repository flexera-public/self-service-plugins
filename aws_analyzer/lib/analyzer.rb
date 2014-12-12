module Analyzer

  # Generic analyzer
  class Analyzer

    # Initialize
    #
    # @option [String] :cloud Service cloud provider, 'aws' or 'gce'
    def initialize(options)
      cloud = options.delete(:cloud) || 'aws'
      klass = ::Analyzer.const_get(cloud.upcase).const_get('Analyzer') rescue nil
      if klass.nil?
        puts "No analyzer implemented yet for #{cloud.inspect} - exiting..."
        exit 1
      end
      @analyzer = klass.new(options)
    end

    # Analyze service from definition in hash
    def analyze_service(service, force=false)
      begin
        definition = @analyzer.analyze(service)
        errors = @analyzer.errors
      rescue Exception => e
        if e.is_a? SystemExit
          raise
        end
        errors = [e.message + " from\n" + e.backtrace.join("\n")]
      end
      if errors && !errors.empty?
        if errors.size == 1
          puts "ERROR: #{errors.first}"
        else
          puts "ERRORS:\n#{errors.join("\n")}"
        end
        exit 1 if !force
        puts
      end
      puts definition.to_yaml unless definition.nil? || definition.empty?
    end

    # YAML representation
    def to_yaml
      YAML.dump({ name:               @name,
                  shape:              @shape.to_hash,
                  primary_id:         @primary_id,
                  secondary_ids:      @secondary_ids,
                  actions:            @actions.map(&:to_hash),
                  custom_actions:     @custom_actiohns.map(&:to_hash),
                  collection_actions: @collection_actions.map(&:to_hash),
                  links:              @links })
    end
    alias :to_s :to_yaml

  end

end
