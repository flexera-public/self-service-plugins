# Service definition NOT USED AT THIS POINT
class ServiceDefinition

  # [String] Service name
  attr_reader :name

  # [String] Path prefix for all actions
  attr_reader :path_prefix

  # [String] API Version
  attr_reader :version

  # [Hash<String, Hash>] Hash of resources indexed by name
  # Each resource is represented as a hash, the resource hash has the following keys:
  #   - path [String] Path to resource
  #   - description [String] Optional description
  #   - actions [Hash] Hash of actions indexed by name, each action is a hash (described below)
  # Resource actions are also represented as hashes with the following keys:
  #   - verb [String] HTTP verb for action, one of 'get', 'post', 'put' or 'delete'
  #   - path [String] Action path (URL is built using the service path prefix + resource path + action path)
  #   - params [Hash] Hash of action params indexed by name, each param is a hash (described below)
  #   - payload [Hash] Hash of action payload fields indexed by name, each field is a hash (described below)
  #   - headers [Hash] hash of request headers indexed by name, each header is a hash (described below)
  #   - responses [Hash] hash of responses indexed by name, each response is a hash (described below)
  # Resource action responses are hashes with the following keys:
  #   - status [Integer]
  #   - description [String]
  #   - headers [Hash] Hash of headers indexed by name, value is string
  #   - media_type [Hash] Hash representing media type
  # Params, payload fields and resource action headers hashes:
  #   - type [String], one of 'string', 'integer'
  attr_reader :resources

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

