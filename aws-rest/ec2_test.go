// Copyright (c) 2015 RightScale, Inc. - see LICENSE

package main

// Omega: Alt+937

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EC2 service", func() {

	client := &http.Client{}

	It("lists zones", func() {
		resp, err := client.Get(SRV + "/rest/ec2/us-east-1/availability_zones")
		Ω(err).ShouldNot(HaveOccurred())
		body, err := ioutil.ReadAll(resp.Body)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(body).Should(ContainSubstring("zone_name"))
		Ω(resp.StatusCode).Should(Equal(200))
	})

	/*
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
	    resp = delete '/ec2/key_pairs/deleteme_now', Yajl::Encoder.encode(args), "CONTENT_TYPE" => "application/json"
	    put_response(resp)

	    args = { key_name: "deleteme_now" }
	    resp = post_json '/ec2/key_pairs', args
	    put_response(resp)
	    expect(resp.status).to eq(201)
	    expect(resp.location).to match("deleteme_now")

	    resp = delete '/ec2/key_pairs/deleteme_now'
	    put_response(resp)
	    expect(resp.status).to eq(204)
	  end

	  it 'launches and terminates an instance' do
	    # we start by deleting the instance in case it exists
	    #resp = delete '/ec2/load_balancers/deleteme-now'
	    #put_response(resp)

	    args = { image_id: "ami-018c9568", min_count: 1, max_count: 1 }
	    resp = post_json '/ec2/instances/actions/run', args
	    put_response(resp)
	    expect(resp.status).to eq(200)
	    expect(resp.body).to match("instance_id")
	    r = Yajl::Parser.parse(resp.body)
	    expect(r['instances'].size).to eq(1)
	    instance_id = r['instances'].first["instance_id"]

	    loop do
	      resp = get "/ec2/instances/#{instance_id}"
	      expect(resp.status).to eq(200)
	      r = Yajl::Parser.parse(resp.body)
	      puts "State: #{r['state']['name']}"
	      break if r['state']['name'] == 'booting' || r['state']['name'] == 'running'
	      sleep 5
	    end

	    args = { instance_ids: [ instance_id ] }
	    resp = post_json '/ec2/instances/actions/terminate', args
	    put_response(resp)
	    expect(resp.status).to eq(200)
	  end

	*/

})
