require_relative 'spec_helper.rb'
require 'json'

describe 'RDS' do
  it 'lists instances' do
    resp = get '/rds/db_instances'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("db_instance_identifier")
  end

  it 'shows one instance' do
    # start with a listing
    resp = get '/rds/db_instances'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("db_instance_identifier")

    r = Yajl::Parser.parse(resp.body)
    resp = get "/rds/db_instances/#{r.first['db_instance_identifier']}"
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("db_instance_identifier")
  end

  it 'creates and deletes an instance' do
    # we start by deleting the DB instance in case it already exists
    loop do
      args = { skip_final_snapshot: true }
      resp = delete '/rds/db_instances/deleteme-now', Yajl::Encoder.encode(args),
        "CONTENT_TYPE" => "application/json"
      put_response(resp)
      break if resp.status == 204 || resp.status == 404 || resp.status == 400 && resp.body =~ /already being deleted/
      sleep 5
    end

    # we now wait for it to be gone
    loop do
      resp = get '/rds/db_instances/deleteme-now'
      #put_response(resp)
      break if resp.status == 404
      r = Yajl::Parser.parse(resp.body)
      puts "Status: #{r["db_instance_status"]} (waiting fordeletion to complete)"
      sleep 5
    end

    # let's actually create the DB instance
    args = { db_instance_identifier: "deleteme-now", allocated_storage: 5,
             db_instance_class: 'db.t1.micro', engine: 'MySQL',
             master_username: 'notme', master_user_password: 'abc_def$123' }
    resp = post_json '/rds/db_instances', args
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme-now")

    # now wait for it to become available
    loop do
      resp = get '/rds/db_instances/deleteme-now'
      #put_response(resp)
      r = Yajl::Parser.parse(resp.body)
      puts "Status: #{r["db_instance_status"]} (waiting for 'available')"
      break if resp.status == 200 && resp.body =~ /"available"/
      sleep 5
    end

    # and finally we delete it again (we won't wait, though)
    args = { skip_final_snapshot: true }
    resp = delete '/rds/db_instances/deleteme-now', Yajl::Encoder.encode(args),
      "CONTENT_TYPE" => "application/json"
    put_response(resp)
    expect(resp.status).to eq(204)
  end

end
