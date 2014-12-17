AWS-CONNECTOR
=============

This app provides a uniform JSON HTTP interface to all AWS services. It allows clients to use
a single form of call to perform any AWS service API request. All requests are POSTs, they all
use a URI of the form `/service/call` (for example, `/ec2/run-instances`), they all accespt
a JSON request body, and they all produce a JSON response body.

Structure
---------
This is a sinatra app with currently only a single controller in `app/proxy.rb`. It loads the
aws-sdk gem in order to get access to the AWS API metadata and to perform actual queries in AWS.
When the handler receives a request, it verifies that the requested service and call exist by
inspecting the gem, then converts the receives payload into ruby hashes & arrays, invokes the
gem to perform the API request, and finally converts the returned result ruby objects back into
JSON for the response.

How-To
------
Ultimetaly this app should run under rainbows. In the meantime it can be run using rackup.
In order to try something out there are test specs for a couple of services, use something like
`bundle exec rspec spec/ec2.rb`. In order this to work you must have your AWS creds in the environment
(`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`).
