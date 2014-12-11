#! /bin/bash -e

echo "{" > aws.json
#for f in $1/*.api.json; do
for f in $1/{EC2,ElasticLoadBalancing,RDS,CloudFormation}.api.json; do
  name=`basename $f .api.json`
  echo "\"${name}\":" >>aws.json
  cat $f >>aws.json
  echo "," >>aws.json
done
echo "}" >>aws.json
