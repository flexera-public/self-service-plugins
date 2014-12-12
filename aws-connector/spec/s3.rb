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

describe 'S3' do
  it 'lists buckets' do
    args = { }
    resp = post_json '/s3/list_buckets', args
    put_response(resp)
    expect(resp.status).to eq(200)
    expect(resp.body).to match("creation_date")
  end

  it 'lists a big bucket' do
    args = { bucket: 'test-us3-cloud-trail' }
    resp = post_json '/s3/list_objects', args
    #put_response(resp)
    expect(resp.status).to eq(200)
    data = Yajl::Parser.parse(resp.body)
    expect(data.keys.first).to eq('contents')
    expect(data['contents'].size).to be > 4000
  end
end
