require 'uri'
require 'net/http'
require 'json'

  def get(path)
    request(path) do |uri|
      req = Net::HTTP::Get.new(uri.path)
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

  get('http://localhost:8082/ec2cft/accounts/60073/stacks')


