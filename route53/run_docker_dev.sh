#!/bin/sh -e
docker run --rm -it -v $(pwd):/opt/praxis -p 8888:8888 -e "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" -e "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" -e "API_SHARED_SECRET=$API_SHARED_SECRET" rgeyer/route53praxis bundle exec rainbows
