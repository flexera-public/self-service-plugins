require_relative 'spec_helper.rb'
require 'json'

describe 'ELB' do
  it 'lists load balancers' do
    args = { }
    resp = get '/elastic_load_balancing/load_balancers'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("load_balancer_name")
  end

  it 'shows a load balancer' do
    args = { }
    resp = get '/elastic_load_balancing/load_balancers/deleteme-now'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("deleteme-now")
  end

  it 'creates and deletes a load balancer' do
    # we start by deleting the load balancer in case it exists
    resp = delete '/elastic_load_balancing/load_balancers/deleteme_now'
    put_response(resp)

    args = { load_balancer_name: "deleteme_now" }
    resp = post_json '/elastic_load_balancing/load_balancers', args
    put_response(resp)
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme_now")

    resp = delete '/elastic_load_balancing/load_balancers/deleteme_now'
    put_response(resp)
    expect(resp.status).to eq(204)
  end

end
