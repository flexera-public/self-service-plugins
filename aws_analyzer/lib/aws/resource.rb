module Analyzer

  module AWS

    # A service resource, the main point of this class is to make sure that we can easily identify operations that apply
    # to either the resourcd or the collection.
    class Resource

      attr_reader :name

      # Initialize with resource name
      def initialize(name)
        @name = name.underscore
        @actions = {}
        @collection_actions = {}
        @custom_actions = {}
      end

      # Register operation
      def add_operation(name, res_name, op, is_collection)
        truncate_size = res_name.size + 1
        n = name[0..-truncate_size].underscore
        if n == 'describe'
          n = is_collection ? 'index' : 'show'
        end
        operation = to_operation(op, n)
        if is_collection
          if n == 'index'
            @actions['index'] = operation
          else
            @collection_actions[n] = operation
          end
        else
          if n == 'show'
            # Let's set the shape of the resource with the result of a describe
            shape = op['output']['shape']
            if shape.nil?
              raise "No shape for describe??? Resource: #{name}, Operation: #{op['name']}"
            end
          end
          if ['create', 'delete', 'update', 'show'].include?(n)
            @actions[n] = operation
          else
            @custom_actions[n] = operation
          end
        end
      end

      # Map raw JSON operation to analyzed YAML operation
      # e.g.
      #    - name: DescribeStackResource
      #      http:
      #        method: POST
      #        requestUri: "/"
      #      input:
      #        shape: DescribeStackResourceInput
      #      output:
      #        shape: DescribeStackResourceOutput
      #        resultWrapper: DescribeStackResourceResult
      # becomes
      #    - name: show
      #      verb: post
      #      path: "/"
      #      payload: describe_stack_resource_input
      #      params:
      #      response: describe_stack_resource_output
      def to_operation(op, name)
        { name:     name,
          verb:     op['http']['method'].downcase,
          path:     op['http']['requestUri'],
          payload:  op['input']['shape'].underscore,
          params:   [],
          response: (out = op['output']) && out['shape'].underscore }
      end

      # Hashify
      def to_hash
        { name:               @name,
          shape:              @shape,
          primary_id:         @primary_id,
          secondary_ids:      @secondary_ids,
          actions:            @actions.values,
          custom_actions:     @custom_actions.values,
          collection_actions: @collection_actions.values }
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
      # Create resource if non-existent, checks whether operation is collection or resource operation
      def add_resource_operation(res_name, op_name, op)
        singular = res_name.singularize
        collection_operation = (singular != res_name)
        res = @resources[singular] ||= Resource.new(singular)
        res.add_operation(op_name, res_name, op, collection_operation)
      end

      # Known resource names
      def resource_names
        @resources.keys
      end

    end

  end

end
