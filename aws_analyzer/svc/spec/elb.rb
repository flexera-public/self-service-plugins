require_relative 'spec_helper.rb'
require 'json'

describe 'ELB' do
  it 'lists load balancers' do
    resp = get '/elastic_load_balancing/load_balancers'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("load_balancer_name")
  end

  it 'shows a load balancer' do
    resp = get '/elastic_load_balancing/load_balancers'
    r = Yajl::Parser.parse(resp.body).first
    resp = get "/elastic_load_balancing/load_balancers/#{r['load_balancer_name']}"
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("load_balancer_name")
  end

  it 'creates and deletes a load balancer' do
    # we start by deleting the load balancer in case it exists
    resp = delete '/elastic_load_balancing/load_balancers/deleteme-now'
    put_response(resp)

    args = { load_balancer_name: "deleteme-now", listeners: [
      protocol: 'http', load_balancer_port: 80, instance_protocol: 'http', instance_port: 80 ],
      availability_zones: [ 'us-east-1b' ] }
    resp = post_json '/elastic_load_balancing/load_balancers', args
    put_response(resp)
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme-now")

    resp = delete '/elastic_load_balancing/load_balancers/deleteme-now'
    put_response(resp)
    expect(resp.status).to eq(204)
  end

end
