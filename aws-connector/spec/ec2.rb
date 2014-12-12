require_relative 'spec_helper.rb'

def post_json(uri, args)
  post uri, Yajl::Encoder.encode(args), "CONTENT_TYPE" => "application/json"
end

def put_response(resp)
  if resp.status < 300
    puts "OK: #{resp.body}"
  else
    puts "ERROR #{resp.status}: #{resp.body}"
  end
end

describe 'EC2' do
  it 'lists zones' do
    args = { }
    resp = post_json '/ec2/describe_availability_zones', args
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("zone_name")
  end

  it 'returns argument errors' do
    args = { stack_name: "teststack" }
    resp = post_json '/ec2/describe_availability_zones', args
    put_response(resp)
    expect(resp.status).to eq(400)
    expect(resp.body).to match("stack_name")
  end

  it 'allocates and deallocates an EIP' do
    args = { domain: 'vpc' }
    resp = post_json '/ec2/allocate_address', args
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("allocation_id")

    aid = Yajl::Parser.parse(resp.body)['allocation_id']
    args = { allocation_id: aid }
    resp = post_json '/ec2/release_address', args
    put_response(resp)
    expect(resp.status).to eq(200)
  end
end
