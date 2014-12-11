require 'json'

# Resource operation prefixes
RESOURCE_ACTIONS = ['Describe', 'Create', 'Update', 'Modify', 'Change', 'Delete', 'Get', 'List', 'Put']

# Infamous snakecase
class String
  def underscore
    self.gsub(/::/, '/').
    gsub(/([A-Z]+)([A-Z][a-z])/,'\1_\2').
    gsub(/([a-z\d])([A-Z])/,'\1_\2').
    tr("-", "_").
    downcase
  end
end

# AWS Service definition: List of resources and independent operations (that need to be mapped to resource "manually")
class AWSServiceDefinition

  # [String] Service name
  attr_reader :name

  # [String] Endpoint prefix
  attr_reader :endpoint

  # [Hash] Metadata hash as defined in JSON
  attr_reader :metadata

  # [Hash<String, ServiceResource>] Hash of resources indexed by name
  # Each resource is represented as a hash of operations indexed by operation name
  attr_reader :resources

  # [Hash<String, Hash>] Hash of operations indexed by name
  attr_reader :operations

  # [Hash<String, Hash>] Hash of shapes indexed by name
  attr_reader :shapes

  def initialize(options)
    options.each { |k ,v| instance_variable_set("@#{k}", v) }
  end

  # Friendly human representation
  def to_s
    res = ["==== #{@name} ===="]
    if @resources.size > 0
      res << '== Resources'
      res += @resources.map do |n, r|
        resource_actions = r.resource_operations.keys.map(&:underscore).join(', ')
        if resource_actions.size > 0
          resource_actions = " - resource: #{resource_actions}"
        end
        collection_actions = r.collection_operations.keys.map(&:underscore).join(', ')
        if collection_actions.size > 0
          collection_actions = " - collection: #{collection_actions}"
        end
        "  #{n}#{resource_actions}#{collection_actions}"
      end
    end
    if @operations.size > 0
      res << '== Unidentified Operations'
      res += @operations.keys.map { |k| "  #{k}" }
    end
    res.join("\n")
    res << "\n"
  end

end

# A service resource, the main point of this class is to make sure that we can easily identify operations that apply
# to either the resourcd or the collection.
class ServiceResource

  attr_reader :name
  attr_reader :resource_operations
  attr_reader :collection_operations

  # Initialize with resource name
  def initialize(name)
    @name = name.chomp('s')
    @resource_operations = {}
    @collection_operations = {}
  end

  # Register operation
  def add_operation(name, op, is_collection)
    truncate_size = @name.size + 1 + (is_collection ? 1 : 0)
    n = name[0..-truncate_size]
    if is_collection
      @collection_operations[n] = op
    else
      @resource_operations[n] = op
    end
  end

end

# Registry of resources for a given service
class ServiceResourceRegistry

  # Resources indexed by name
  attr_reader :resources

  def initialize
    @resources = {}
  end

  # Add operation to resource
  # Create resource if non-existent, checks whether operation is collection or resource operation
  def add_resource_operation(res_name, op_name, op)
    collection_operation = (res_name[-1] == 's')
    singular = res_name.chomp('s')
    res = @resources[singular] ||= ServiceResource.new(singular)
    res.add_operation(op_name, op, collection_operation)
  end

  # Known resource names
  def resource_names
    @resources.keys
  end

end

# Anayze a single service hash
class ServiceAnalyzer

  # Analyze 'operations' hash for given service
  # Returns a ServiceDefinition object (contains list of resources and custom actions)
  def analyze(service)
    registry = ServiceResourceRegistry.new
    operations = service['operations'].keys
    candidates, remaining = operations.partition { |o| is_resource_action?(o) }

    # 1. Identify resources by using well known operation prefixes
    candidates.each do |c|
      name = c.gsub(/^(#{RESOURCE_ACTIONS.join('|')})/, '')
      registry.add_resource_operation(name, c, service['operations'][c])
    end

    # 2. Collect custom actions that apply to identified resources
    matched = []
    remaining.each do |r|
      candidates = registry.resource_names.select { |n| r =~ /(#{n}|#{n}s)$/ }
      next if candidates.empty?
      matched << r
      candidates.sort { |a, b| b.size <=> a.size } # Longer name first
      name = candidates.first
      registry.add_resource_operation(name, r, service['operations'][r])
    end

    # 3. Collect remaining - unidentified operations
    operations = (remaining - matched).inject({}) { |m, o| m[o] = service['operations'][o]; m }

    AWSServiceDefinition.new(name:       service['metadata']['serviceFullName'],
                          endpoint:   service['metadata']['endpointPrefix'],
                          metadata:   service['metadata'],
                          resources:  registry.resources,
                          operations: operations,
                          shapes:     service['shapes'])
  end

  # true if name is an operation on a resource (i.e. has a well-known prefix)
  def is_resource_action?(name)
    !!RESOURCE_ACTIONS.any? { |s| name =~ /^#{s}/ }
  end

end

