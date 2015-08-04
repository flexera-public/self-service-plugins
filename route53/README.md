# Route53 Praxis App
A simple Praxis App to perform Route53 CRUD requests from Self Service

## Limitations
It's single tenant, with the AWS credentials provided as environment variables (see [Credentials][])

### Zones
It doesn't support [private zones](http://docs.aws.amazon.com/Route53/latest/DeveloperGuide/hosted-zones-private.html)

It doesn't support [delegation sets](http://docs.aws.amazon.com/Route53/latest/APIReference/actions-on-reusable-delegation-sets.html)

### Records
It doesn't support any of the custom [routing policies](http://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy.html)

## Credentials
Requires a single set of AWS credentials provided as environment variables.

* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`

It also requires an API shared secret. This shared secret should be known only to your instance of this app, and the RightScale SelfService namespace you use to connect to it. It will be used to establish trust between these parties.

You provide the API shared secret as an environment variable

`API_SHARED_SECRET`
