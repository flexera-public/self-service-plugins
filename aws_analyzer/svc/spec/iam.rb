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

  it 'shows a user' do
    args = { }
    resp = get '/iam/users/raphael'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("user_name")
  end

  it 'finds a group' do
    # this doesn't work...
    resp = get '/iam/groups?filter[]=path_prefix==/'
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("power-users")
  end

  it 'creates and deletes a user' do
    # we start by deleting the user in case it exists
    args = { user_name: "deleteme_now" }
    resp = post_json '/iam/groups/power-users/actions/remove_user_from', args
    put_response(resp)
    resp = delete '/iam/users/deleteme_now'
    put_response(resp)

    args = { user_name: "deleteme_now" }
    resp = post_json '/iam/users', args
    put_response(resp)
    expect(resp.status).to eq(201)
    expect(resp.location).to match("deleteme_now")

    args = { user_name: "deleteme_now" }
    resp = post_json '/iam/groups/power-users/actions/add_user_to', args
    put_response(resp)
    expect(resp.status).to eq(204)

    args = { user_name: "deleteme_now" }
    resp = post_json '/iam/groups/power-users/actions/remove_user_from', args
    put_response(resp)
    expect(resp.status).to eq(204)

    resp = delete '/iam/users/deleteme_now'
    put_response(resp)
    expect(resp.status).to eq(204)
  end

end
