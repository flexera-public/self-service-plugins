module V1
  class Stack
    include Praxis::Controller

    implements V1::ApiResources::Stack

    def index(account_id:, **params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      cfm = Aws::CloudFormation::Client.new
      stacks = cfm.describe_stacks.stacks
      prax_stax = Array.new

      stacks.each do |s| 
        next if s["stack_status"] == "DELETE_COMPLETE" 
        stack = {}

        stack["id"] = s["stack_id"]
        stack["href"] = "/ec2cft/accounts/" + account_id.to_s + "/stacks/" + s["stack_name"].to_s
        stack["name"] = s["stack_name"]
        stack["status"] = s["stack_status"]
        stack["status_reason"] = s["stack_status_reason"]
        stack["creation_time"] = s["creation_time"]
        stack["description"] = s["description"]

        prax_stax.push(stack)
      end

      response.headers['Content-Type'] = 'application/json'

      response.body = JSON.pretty_generate(prax_stax.map { |s| V1::MediaTypes::Stack.dump(s) })
      response
    end

    def show(account_id:, name:, view:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      cfm = Aws::CloudFormation::Client.new

      begin
        resp = cfm.describe_stacks({
          stack_name: name,
        })

        stack = resp.stacks[0]

        outputs = Array.new
        stack.outputs.each do |o|
          op = { 
            "key" => o.output_key,
            "value" => o.output_value,
            "description" => o.description
          }
          outputs.push(op)
        end
        resp = {
          "id" => stack.stack_id,
          "href" => "/ec2cft/accounts/" + account_id.to_s + "/stacks/#{stack.stack_name}",
          "creation_time" => stack.creation_time,
          "description" => stack.description,
          "name" => stack.stack_name,
          "status" => stack.stack_status,
          "status_reason" => stack.stack_status_reason,
          "outputs" => outputs
        }  
        response.body = JSON.pretty_generate(V1::MediaTypes::Stack.dump(resp, :view=>(view ||= :default).to_sym))
      rescue Exception => e  
        self.response = Praxis::Responses::NotFound.new()
        response.body = { error: '404: Not found' }
      end
      response.headers['Content-Type'] = 'application/json'
      response
    end

    def create(account_id:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      cfm = Aws::CloudFormation::Client.new

      params = []
      request.payload.parameters.each do |k,v|
        params << {parameter_key: k, parameter_value: v}
      end

      begin
        # Create the stack
        stack_id = cfm.create_stack({
          stack_name: request.payload.name, 
          template_url: request.payload.template, 
          parameters: params
        }).stack_id

        # Get the stack properties
        resp = cfm.describe_stacks({
          stack_name: stack_id,
        })
        stack = resp.stacks[0]

        # Form the response
        self.response = Praxis::Responses::Created.new()
        resp = {
          "id" => stack.stack_id,
          "href" => "/ec2cft/accounts/" + account_id.to_s + "/stacks/#{stack.stack_name}",
          "creation_time" => stack.creation_time,
          "description" => stack.description,
          "name" => stack.stack_name,
          "status" => stack.stack_status,
          "status_reason" => stack.stack_status_reason
        }  
        response.headers['Location'] = resp["href"]
        response.body = JSON.pretty_generate(V1::MediaTypes::Stack.dump(resp))        
      rescue Exception => e  
        self.response = Praxis::Responses::UnprocessableEntity.new()
        response.body = { error: '422: Not able to create record: ' + e.message }
      end
      response.headers['Content-Type'] = 'application/json'      
      response

    end

    def delete(name:, **other_params)
      resp = authenticate!(request.headers["X_Api_Shared_Secret"])
      return resp if resp

      cfm = Aws::CloudFormation::Client.new

      resp = cfm.delete_stack({
        stack_name: name
      })

      # AWS appears to always return 'successful', even if the stack name doesn't exist
      if resp.successful?
        self.response = Praxis::Responses::NoContent.new()   
      else
        self.response = Praxis::Responses::UnprocessableEntity.new()
        response.headers['Content-Type'] = 'application/json'      
        response.body = { error: '422: Not able to delete stack with name: ' + name + '. Response was: ' + resp.data.to_s }        
      end

      response

    end

    def update(domain:, id:, **other_params)

      # api = get_api
      # res = api.update_record(domain, id, request.payload.name, request.payload.type, request.payload.value)

      # if res["error"].nil? 
      #   self.response = Praxis::Responses::NoContent.new()
      #   response.body = res
      # else
      #   self.response = Praxis::Responses::UnprocessableEntity.new()
      #   response.headers['Content-Type'] = 'application/json'      
      #   response.body = { error: '422: Not able to update record' }
      # end
      # response

    end 

    private

    def authenticate!(secret)
      if secret != ENV["PLUGIN_SHARED_SECRET"]
        self.response = Praxis::Responses::Forbidden.new()
        response.body = { error: '403: Invalid shared secret'}
        return response
      else
        return nil
      end
    end

  end
end
