#!/bin/sh -e
docker run --rm -it -v $(pwd):/opt/praxis -p 8888:3000 -e "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" -e "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" rgeyer/route53praxis bundle exec thin start
