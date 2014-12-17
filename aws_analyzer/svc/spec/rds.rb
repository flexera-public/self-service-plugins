require_relative 'spec_helper.rb'
require 'json'

describe 'RDS' do
  it 'lists instances' do
    args = { }
    resp = get '/rds/db_instances'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("user_name")
  end

  it 'creates and deletes an instance' do
    # we start by deleting the user in case it exists
    resp = delete '/rds/db_instances/deleteme_now'
    put_response(resp)

    args = { db_instance_identifier: "deleteme_now", allocated_storage: 5,
             db_instance_class: 'db.t1.micro', engine: 'MySQL',
             master_username: 'notme', master_user_password: '@&*$%$@^3237743676' }
    resp = post_json '/rds/db_instances', args
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme_now")
  end

end
