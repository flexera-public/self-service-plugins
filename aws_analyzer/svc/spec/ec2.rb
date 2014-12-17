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
    resp = delete '/ec2/key_pair/deleteme_now', Yajl::Encoder.encode(args), "CONTENT_TYPE" => "application/json"
    put_response(resp)

    args = { key_name: "deleteme_now" }
    resp = post_json '/ec2/key_pair', args
    put_response(resp)
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme_now")

    resp = delete '/ec2/key_pair/deleteme_now'
    put_response(resp)
    expect(resp.status).to eq(204)
  end

=begin
  it 'allocates and deallocates an EIP' do
    args = { domain: 'vpc' }
    resp = post_json '/ec2/address', args
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("allocation_id")

    #aid = Yajl::Parser.parse(resp.body)['allocation_id']
    #args = { allocation_id: aid }
    #resp = post_json '/ec2/release_address', args
    #put_response(resp)
    #expect(resp.status).to eq(200)
  end
=end


=begin
  it 'returns argument errors' do
    args = { stack_name: "teststack" }
    resp = post_json '/ec2/describe_availability_zones', args
    put_response(resp)
    expect(resp.status).to eq(400)
    expect(resp.body).to match("stack_name")
  end

=end
end
