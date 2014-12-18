module Analyzer

  module AWS

    # Utility class used to convert operations coming from service json files to action
    class Operation

      # Map raw JSON operation to analyzed YAML action
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
      def self.to_action(op)
        ::Analyzer::Action.new(name:      op['name'].underscore,
                               verb:      op['http']['method'].downcase,
                               path:      op['http']['requestUri'],
                               payload:   (sh = op['input'] && op['input']['shape']) ? sh.underscore : op['input'],
                               params:    [],
                               response:  (out = op['output']) && out['shape'].underscore)
      end

    end

  end

end


