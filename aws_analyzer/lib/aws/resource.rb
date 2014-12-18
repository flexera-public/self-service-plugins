module Analyzer

  module AWS

    # Resources that have a plural name
    PLURAL_RESOURCE_NAMES = ['DhcpOptions']

    # A service resource, the main point of this class is to make sure that we can easily identify operations that apply
    # to either the resourcd or the collection.
    class Resource

      # [String] Resource name (e.g. "stack")
      attr_reader :name

      # [String] Resource original name (e.g. "Stack")
      attr_reader :orig_name

      # [String] Resource shape name (e.g. "Stack")
      attr_reader :shape

      # [String] Resource primary id field (e.g. "StackId")
      # TBD
      attr_accessor :primary_id

      # [Array<String>] Resource secondary ids field (e.g. ["StackName"])
      # TBD
      attr_accessor :secondary_ids

      # [Hash<String, Action>] Resource CRUD actions (index, show, update, create, delete)
      attr_reader :actions

      # [Hash<String, Action>] Resource custom actions (e.g. cancel_update)
      attr_reader :custom_actions

      # [Hash<String, Action>] Resource collection custom actions (e.g. list)
      attr_reader :collection_actions

      # [Hash<String, String>] Linked resource names indexed by link field name (e.g. { "stack_id" => "Stack" })
      # TBD
      attr_reader :links

      # Initialize with resource name
      def initialize(name)
        @name               = name.underscore
        @orig_name          = name
        @actions            = {}
        @collection_actions = {}
        @custom_actions     = {}
      end

      # Register operation
      # OK, here is the trick:
      # name is the CamelCase name of the operation, this name ends with either the ResourceName or ResourceNames
      # for operations that apply to the collection. We detect which one it is and then infer the final action
      # name and type (resource, collection or custom) from that.
      def add_operation(op, shapes)
        name = op['name']
        is_collection = name !~ /#{@orig_name}$/ # @orig_name is the singular version of ResourceName
        n = name.gsub(/(#{@orig_name}|#{@orig_name.pluralize})$/, '').underscore
        if n == 'describe' || n == 'list' || n == 'get'
          n = is_collection ? 'index' : 'show'
        end
        action = Operation.to_action(op)
        if is_collection
          if n == 'index'
            @actions['index'] = action
            # Some resource only have an index action - no show action (index can be filtered by id)
            # So try to infer shape from index if not set yet
            if @shape.nil?
              oshape = (os = op['output']['shape']) && shapes[os]
              unless oshape.nil?
                oshape['members'].each do |sn, m|
                  if name =~ /#{sn.singularize}/ # e.g. 'DescribeStackEvents' =~ 'StackEvents'
                    @shape = shapes[m['shape']]['member']['shape'].underscore rescue nil
                  end
                end
              end
            end
          else
            @collection_actions[n] = action
          end
        else
          if n == 'show'
            # Let's set the shape of the resource with the result of a describe
            candidate = op['output']['shape']
            if candidate.nil?
              raise "No shape for describe??? Resource: #{name}, Operation: #{op['name']}"
            end
            cs = shapes[candidate]
            shape_member = (smn = cs['members'].keys.detect { |k| op['name'] =~ /#{k}/ }) && cs['members'][smn]
            if shape_member.is_a?(Hash) && shape_member.keys.first == 'shape'
              @shape = shape_member['shape'].underscore
            else
              @shape = candidate.underscore
            end
          end
          if ['create', 'delete', 'update', 'show'].include?(n)
            @actions[n] = action
          else
            @custom_actions[n] = action
          end
        end
      end


      # Hashify
      def to_hash
        { 'name'               => @name,
          'shape'              => @shape,
          'primary_id'         => @primary_id,
          'secondary_ids'      => @secondary_ids,
          'actions'            => @actions.inject({}) { |m, (k, v)| m[k] = v.to_hash; m },
          'custom_actions'     => @custom_actions.inject({}) { |m, (k, v)| m[k] = v.to_hash; m },
          'collection_actions' => @collection_actions.inject({}) { |m, (k, v)| m[k] = v.to_hash; m } }
      end

    end

    # Registry of resources for a given service
    class ResourceRegistry

      # Resources indexed by name
      attr_reader :resources

      def initialize
        @resources = {}
      end

      # Add operation to resource
      # Create resource if non-existent, checks whether operation is collection or resource action
      # Takes shapes define in API json so that it can apply a heuristic to infer the name of the resource shape
      # for resources that only support an index call (and no show call)
      def add_resource_operation(res_name, op, shapes)
        canonical = canonical_name(res_name)
        res = @resources[canonical] ||= Resource.new(canonical)
        res.add_operation(op, shapes)
      end

      # Known resource names
      def resource_names
        @resources.keys
      end

      # Singularize name unless it's in the exception list
      def canonical_name(base_name)
        PLURAL_RESOURCE_NAMES.include?(base_name) ? base_name : base_name.singularize
      end

      # Remove resources that couldn't be completly identified
      # Returns list of deleted resource names
      def delete_incomplete_resources
        deleted = []
        @resources.values.each do |r|
          if r.primary_id.nil? || r.shape.nil?
            deleted << r.name
          end
        end
        @resources.delete_if { |n, _| deleted.include?(n) }
        deleted
      end

    end

  end

end
