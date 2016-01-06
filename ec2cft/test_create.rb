require 'uri'
require 'net/http'
require 'json'

  def post(path, body)
    request(path) do |uri|
      req = Net::HTTP::Post.new(uri.path)
      req.body = body.to_json
      req
    end

  end

  def request(path)
    uri = URI("#{path}")

    http = Net::HTTP.new(uri.host, uri.port)

    request = yield(uri)

    request_headers.each do |key, value|
      request[key] = value
    end

    response = http.request(request)
    puts response.body
    unparsed_json = response.body == "" ? "{}" : response.body

    puts response.code
    puts response.message
    puts response.header.to_hash

    JSON.parse(unparsed_json)
  end

  def request_headers
    {
      'Content-type' => 'application/json',
      'X-Api-Version' => '1.0',
      "X-Api-Shared-Secret" => "blah"      
    }
  end

  # post('http://localhost:8082/ec2cft/accounts/60073/stacks', {"name"=>"ROLtestagain1","template"=>"https://s3.amazonaws.com/rol-cf-templates/WordPress_Single_Instance_SS.template","parameters"=>[{parameter_key:"DBRootPassword",parameter_value:"blahblah123blahAD"},{parameter_key:"KeyName",parameter_value:"default"}]})
  post('http://localhost:8082/ec2cft/accounts/60073/stacks', {"name"=>"ROLtestinstance","template"=>"https://s3-us-west-2.amazonaws.com/cloudformation-templates-us-west-2/EC2InstanceWithSecurityGroupSample.template","parameters"=>{"KeyName"=>"default","InstanceType"=>"m1.small"}})
  # post('http://localhost:9292/ec2cft/accounts/60073/stacks', {"name"=>"Test5","template"=>"https://s3.amazonaws.com/cloudformation-templates-us-east-1/WordPress_Single_Instance.template","parameters"=>"{\"DBRootPassword\":\"blahblah\",\"KeyName\":\"my normal key\"}"})
  # post('http://54.227.94.207:80/ec2cft/accounts/60073/stacks', {"name"=>"Test4","template"=>"https://s3-external-1.amazonaws.com/cloudformation-samples-us-east-1/Rails_Simple.template","parameters"=>"{\"DBRootPassword\":\"blahblah\",\"KeyName\":\"my normal key\"}"})


