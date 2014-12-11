require_relative '../spec_helper.rb'

describe 'create_instance' do
  it 'creates an instance' do
    post '/ec2/create_instance', {}
    expect(last_response.status).to eq(201)
  end
end
