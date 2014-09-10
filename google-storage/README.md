Google Cloud SQL
================

RightScale Self-Service Plugin to provision Google Cloud SQL databases.

This plugin is a proxy web service that provides a namespace for Google Cloud SQL
allowing databases to be provisioned from CATs. This proxy is multi-tenant, which means
that multiple RightScale accounts can use it at the same time. The mapping from RS account
to Google projects is done in the proxy.

Starting the Proxy
------------------

- Identify or create a Google API project that has the Cloud SQL Admin API enabled
- Create an "Installed application client ID" in a Google project of your choosing
  (has no relationship with the projects where your databases will be)
- Download the client JSON secrets and put them into a file called client_secrets.json
- Install the ruby bundle by running `bundle install`, be sure to use Ruby 1.2.1
- Start the praxis proxy server using `rackup`

Authorizing
-----------

The purpose of this step is to authorize the proxy to access a specific project in Google
Cloud and to establish an association with a RightScale account.

- Identify the RightScale account number which you will be using
- Point your browser at the authentication page:
  `http://localhost:9292/acct/<acct>/auth?project=<proj>`, where <acct> is your RS account
  number and <proj> is your google project name (such as "rightscale.com:test");
  if you can't point your browser at the proxy then use curl and follow the instructions.
- After accepting the access request on Google's site you will be redirected to the local
  proxy with an authentication code.
  If you need to use curl, grab the code from your browser's address bar and complete the
  request with curl to look something like the following:
  `curl http://localhost:9292/acct<acct>/auth/redirect?project=<proj>&code=4/XYAAFHGDGKANBDHSU26487GGJGJH`
- The result of all this is a file `.gc_auth/<acct>` that contains the auth token for this
  account, and the project name.

Caveat: technically you're authorizing the proxy to access *any* project of yours.

Using the proxy
---------------

- List DB instances: `http://localhost:9292/acct/<acct>/instances`

Controllers
-----------

- auth: authentication requests, backed by `lib/auth.rb`
- instances: google cloud sql instances, aka, databases
- hello: the std hello controller installed by the praxis generator, for foolin' around
