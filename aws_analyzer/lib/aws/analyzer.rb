require 'json'

module Analyzer

  module AWS

    # Resource operation prefixes used by heuristic
    RESOURCE_ACTIONS = ['Describe']
    SECONDARY_RESOURCE_ACTIONS = ['List', 'Get']
    ALL_ACTIONS = RESOURCE_ACTIONS + SECONDARY_RESOURCE_ACTIONS

    # Anayze a single service hash
    class Analyzer

      # [Array<String>] Analysis errors
      attr_reader :errors

      # Initialize analyzer with options containing path to JSON definitions
      def initialize(options)
        if (@json_paths = options[:paths]).nil?
          puts 'Please specify path to JSON files with --paths'
          exit 1
        end
      end

      # Analyze 'operations' hash for given service
      # Returns a ServiceDefinition object (contains list of resources and custom actions)
      def analyze(service_name)
        # Locate API and resource json files
        json = ''
        unless @json_paths.detect { |path| json = File.join(path, service_name.camel_case + '.api.json'); File.exist?(json) }
          puts "Hmm there doesn't seem to be a #{service_name.camel_case}.api.json file at #{@json_paths.join(' ,')}, you sure you got that right? (use --paths to specify the location of the JSON files)"
          exit 1
        end

        # read resource definitions
        res_json = []
        @json_paths.each do |path|
          f = File.join(path, service_name.camel_case + '.resources.json')
          res_json << f if File.exist?(f)
        end
        if res_json == []
          puts("Hmm there doesn't seem to be a #{service_name.camel_case}.resources.json file " +
            "at #{@json_paths.join(' ,')}, are you sure you got that right? " +
            "(use --paths to specify the location of the JSON files)")
          exit 1
        end
        # actually load the definitions and merge the top-level values
        resources = {}
        res_json.each do |rj|
          rr = JSON.load(IO.read(rj))
          rr.each_pair do |k,v|
            if resources.key?(k)
              if resources[k].is_a?(Array)
                resources[k] += v
              else
                resources[k].merge!(v)
              end
            else
              resources[k] = v
            end
          end
        end

        service = JSON.load(IO.read(json))
        registry = ResourceRegistry.new
        operations = service['operations'].keys
        # First try with RESOURCE_ACTIONS operation prefixes
        candidates, remaining = operations.partition { |o| is_resource_action?(o, primary_only: true) }
        # If no results then try with SECONDARY_RESOURCE_ACTIONS operation prefixes
        if candidates.empty?
          candidates, remaining = operations.partition { |o| is_resource_action?(o) }
        end

        # 1. Identify resources by using well known operation prefixes
        candidates.each do |c|
          name = c.gsub(/^(#{ALL_ACTIONS.join('|')})/, '')
          if name.underscore =~ /_for_/
            remaining << c
          else
            registry.add_resource_operation(name, service['operations'][c], service['shapes'])
          end
        end

        # 2. Collect custom actions that apply to identified resources
        matched = []
        remaining.each do |r|
          candidates = registry.resource_names.select { |n| r =~ /(#{n}|#{n.pluralize})$/ }
          next if candidates.empty?
          matched << r
          candidates.sort { |a, b| b.size <=> a.size } # Longer name first
          name = candidates.first
          registry.add_resource_operation(name, service['operations'][r], service['shapes'])
        end

        # 3. Leverage resource.json if present to identify id fields and links
        shapes = to_underscore(service['shapes'])
        resources['resources'].each do |name, res|
          next unless res['shape']
          # Let's see if we find an exact match
          existing = registry.resources.select { |n, r| r.orig_name == name }
          if !existing.empty?
            # We found an exact match, use that shape
            if existing.size > 1
              puts "Found ambiguous match: multiple resources with name #{name}..."
              next
            end
            existing.values.first.shape = res['shape'].underscore
          else
            existing = registry.resources.select { |n, r| r.shape == res['shape'].underscore }
            if existing.size > 1
              puts "Found ambiguous match: multiple resources with shape #{res['shape']}..."
              next
            end
          end
          next if existing.empty?
          r = existing.values.first
          shape = shapes[r.shape]
          if shape
            members = shape['members'].keys || []
            res['identifiers'].each do |i|
              candidate = i['memberName'] || i['name']
              next if candidate.nil?
              candidate = candidate.underscore
              if !members.include?(candidate)
                candidate = "#{r.name}_#{candidate}"
                next unless members.include?(candidate)
              end
              if r.primary_id.nil?
                r.primary_id = candidate
              else
                r.secondary_ids ||= []
                r.secondary_ids << candidate
              end
            end
          end
        end

        # 4. Build service level custom actions for unmapped actions
        not_mapped = remaining - matched
        service_actions = not_mapped.inject({}) { |m, nm| m[nm.underscore] = Operation.to_action(service['operations'][nm]); m }

        # 5. Cleanup resources for which we couldn't find an id field
        no_ids = registry.delete_incomplete_resources.sort
        @errors = []
        @errors << "** Failed to find an id for the following resources:\n  #{no_ids.join("\n  ")}" unless no_ids.empty?

        # We're done!
        ::Analyzer::ServiceDefinition.new('name'      => service['metadata']['serviceFullName'],
                                          'url'       => "/aws/#{service['metadata']['endpointPrefix']}",
                                          'metadata'  => service['metadata'],
                                          'resources' => registry.resources.inject({}) { |m, (k, v)| m[k.underscore] = v; m },
                                          'actions'   => service_actions,
                                          'shapes'    => shapes)
      end

      # true if name is an operation on a resource (i.e. has a well-known prefix)
      def is_resource_action?(name, opts={})
        res = !!RESOURCE_ACTIONS.any? { |s| name =~ /^#{s}/ }
        return res if res
        return false if opts[:primary_only]
        res = !!SECONDARY_RESOURCE_ACTIONS.any? { |s| name =~ /^#{s}/ }
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
