require_relative 'spec_helper.rb'
require 'json'

describe 'IAM' do
  it 'lists users' do
    args = { }
    resp = get '/iam/users'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("user_name")
  end

  it 'finds a group' do
    resp = get '/iam/groups?filter[]=path_prefix=power-users'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("power-users")
    expect(resp.body).not_to match("change-password")
  end

  it 'creates and deletes a user' do
    # we start by deleting the user in case it exists
    resp = delete '/iam/users/deleteme_now'
    put_response(resp)

    # find a group to put the user into
    resp = get '/iam/groups?filter[]=path_prefix=power-users'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("power-users")
    expect(resp.body).not_to match("change-password")
    group = Yajl::Parser.parse(resp.body)['group_name']

    args = { user_name: "deleteme_now", group_name: group }
    resp = post_json '/iam/users', args
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme_now")
  end

end
