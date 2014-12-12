require 'json'

module Analyzer

  module AWS

    # Resource operation prefixes used by heuristic
    RESOURCE_ACTIONS = ['Describe']

    # Anayze a single service hash
    class Analyzer

      # [Array<String>] Analysis errors
      attr_reader :errors

      # Initialize analyzer with options containing path to JSON definitions
      def initialize(options)
        if (@json_path = options[:path]).nil?
          puts 'Please specify path to JSON files with --path'
          exit 1
        end
      end

      # Analyze 'operations' hash for given service
      # Returns a ServiceDefinition object (contains list of resources and custom actions)
      def analyze(service_name)
        json = File.join(@json_path, service_name.camel_case + '.api.json')
        if !File.exist?(json)
          puts "Hmm there doesn't seem to be any *.api.json file at #{json}, you sure you got that right? (use --path to specify the location of the JSON files)"
          exit 1
        end
        service = JSON.load(IO.read(json))
        registry = ResourceRegistry.new
        operations = service['operations'].keys
        candidates, remaining = operations.partition { |o| is_resource_action?(o) }

        # 1. Identify resources by using well known operation prefixes
        candidates.each do |c|
          name = c.gsub(/^(#{RESOURCE_ACTIONS.join('|')})/, '')
          registry.add_resource_operation(name, service['operations'][c])
        end

        # 2. Collect custom actions that apply to identified resources
        matched = []
        remaining.each do |r|
          candidates = registry.resource_names.select { |n| r =~ /(#{n}|#{n.pluralize})$/ }
            next if candidates.empty?
          matched << r
          candidates.sort { |a, b| b.size <=> a.size } # Longer name first
          name = candidates.first
          registry.add_resource_operation(name, service['operations'][r])
        end

        # 3. Collect remaining - unidentified operations
        not_mapped = remaining - matched
        @errors = ["Failed to identify a resource for the following operations:\n#{not_mapped.join("\n")}"]

        ::Analyzer::ServiceDefinition.new(name:      service['metadata']['serviceFullName'],
                                          url:       "/aws/#{service['metadata']['endpointPrefix']}",
                                          metadata:  service['metadata'],
                                          resources: registry.resources.inject({}) { |m, (k, v)| m[k.underscore] = v; m },
                                          shapes:    to_underscore(service['shapes']))
      end

      # true if name is an operation on a resource (i.e. has a well-known prefix)
      def is_resource_action?(name)
        !!RESOURCE_ACTIONS.any? { |s| name =~ /^#{s}/ }
      end

      # Keys whose values should not be underscorized - might need to be smarter here
      DO_NOT_UNDERSCORE = ['enum', 'pattern']

      # Recursively traverse data structure and change strings to underscore representation
      def to_underscore(object)
        case object
        when Hash then object.inject({}) do |m, (k, v)|
          if DO_NOT_UNDERSCORE.include?(k)
            m[k.underscore] = v
          else
            m[k.underscore] = to_underscore(v)
          end
          m
        end
        when Array then object.map { |e| to_underscore(e) }
        when String then object.underscore
        else
          object
        end
      end

    end

  end

end
