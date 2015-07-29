Route53 Praxis App
==================
A simple Praxis App to perform Route53 CRUD requests from Self Service

Limitations
-----------
It's single tenant, with the AWS credentials provided as environment variables (see [Credentials][])

It doesn't support [private zones](http://docs.aws.amazon.com/Route53/latest/DeveloperGuide/hosted-zones-private.html)

It doesn't support [delegation sets](http://docs.aws.amazon.com/Route53/latest/APIReference/actions-on-reusable-delegation-sets.html)

Credentials
-----------
Requires a single set of AWS credentials provided as environment variables.

`AWS_ACCESS_KEY_ID`
`AWS_SECRET_ACCESS_KEY`
