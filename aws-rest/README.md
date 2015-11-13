AWS-REST Self-Service and NextGen Plugin
====================================================

Blah blah...

URL paths
---------

This service exposes the following URL roots:
- /aws/* -- performs API requests on AWS services (NextGen connector)
- /rest/* -- exposes a REST interface for all AWS services (SS plugin)
- /xlate/* -- translates from REST to AWS and from AWS to REST (NextGen rest plugin)

The NextGen connector paths are structured as follows:
- POST /aws/:service/:region/:operation
The REST plugin paths are structured as follows:
- GET|POST|PUT|DELETE /rest/:service/:region/
The Xlator paths are structured as follows:
- ?
