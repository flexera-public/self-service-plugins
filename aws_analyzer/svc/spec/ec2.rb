require_relative 'spec_helper.rb'
require 'json'

describe 'EC2' do
  it 'lists zones' do
    args = { }
    resp = get '/ec2/availability_zones'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("zone_name")
  end

  it 'lists key_pairs' do
    args = { }
    resp = get '/ec2/key_pairs'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("key_pairs")
    expect(resp.body).to match("self")
  end

  it 'shows a key_pair' do
    args = { }
    resp = get '/ec2/key_pairs/LLB_SSH1'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("key_name")
    expect(resp.body).to match("self")
  end

  it 'shows a zone' do
    args = { }
    resp = get '/ec2/availability_zones/us-east-1b'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("zone_name")
  end

  it 'allocates and deallocates key pair' do
    args = { key_name: "deleteme_now" }
    resp = delete '/ec2/key_pairs/deleteme_now', Yajl::Encoder.encode(args), "CONTENT_TYPE" => "application/json"
    put_response(resp)

    args = { key_name: "deleteme_now" }
    resp = post_json '/ec2/key_pairs', args
    put_response(resp)
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme_now")

    resp = delete '/ec2/key_pairs/deleteme_now'
    put_response(resp)
    expect(resp.status).to eq(204)
  end

  it 'launches and terminates an instance' do
    # we start by deleting the instance in case it exists
    #resp = delete '/ec2/load_balancers/deleteme-now'
    #put_response(resp)

    args = { image_id: "ami-018c9568", min_count: 1, max_count: 1 }
    resp = post_json '/ec2/instances/actions/run', args
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("instance_id")
    r = Yajl::Parser.parse(resp.body)
    expect(r['instances'].size).to eq(1)
    instance_id = r['instances'].first["instance_id"]

    loop do
      resp = get "/ec2/instances/#{instance_id}"
      put_response(resp)
      expect(resp.status).to eq(200)
      r = Yajl::Parser.parse(resp.body)
      break if r['state'] == 'booting' || r['state'] == 'running'
      sleep 5
    end

    args = { instance_ids: [ instance_id ] }
    resp = post_json '/ec2/instances/actions/terminate', args
    put_response(resp)
    expect(resp.status).to eq(200)
  end

end
