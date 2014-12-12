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

describe 'CloudFormation' do
  it 'lists stacks' do
    args = { stack_name: "teststack" }
    resp = post_json '/cloud_formation/describe_stacks', args
    put_response(resp)
    expect(resp.status).to eq(200)
  end

  it 'handles list stacks error' do
    args = { stack_name: "notexistent" }
    resp = post_json '/cloud_formation/describe_stacks', args
    put_response(resp)
    expect(resp.status).to eq(400)
    expect(resp.body).to eq("Stack:notexistent does not exist")
  end

  it 'gets a template' do
    args = { stack_name: "teststack" }
    resp = post_json '/cloud_formation/get_template', args
    put_response(resp)
    expect(resp.status).to eq(200)
  end

  it 'updates a stack template' do
    args = { stack_name: "teststack" }
    resp = post_json '/cloud_formation/get_template', args
    expect(resp.status).to eq(200)

    template = Yajl::Parser.parse(resp.body)['template_body']
    puts "Template: <<#{template}>>"

    args[:template_body] = template
    resp = post_json '/cloud_formation/update_stack', args
    put_response(resp)
    expect(resp.status).to eq(400)

    template['Description'] = "Description updated at #{Time.now()}"
    args[:template_body] = template
    resp = post_json '/cloud_formation/update_stack', args
    put_response(resp)
    expect(resp.status).to eq(200)
  end

end
