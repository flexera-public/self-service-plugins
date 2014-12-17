AWS-ANALYZE-SVC
===============

This app provides a restful resource JSON HTTP interface to all AWS services. It allows clients to
use a relatively conventional resource interface to access all of AWS.

Structure
---------
This is a sinatra app with currently only a single controller in `app/restifier.rb`.
It loads the analyzer with an in-memory descritpion of the resources in order to figure out
which resources each AWS service has and what the actions are. In order to actually perform
any request it invokes the aws-connector.

How-To
------
Launch the aws-connector on port 9001. Then use something like `bundle exec rspec spec/ec2.rb`
to run a simple test.
