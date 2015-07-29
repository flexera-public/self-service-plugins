FROM rgeyer/praxis
MAINTAINER Ryan J. Geyer <me@ryangeyer.com>

COPY Gemfile /tmp/Gemfile
RUN cd /tmp && bundle install

ENV WORKDIR /opt/praxis
WORKDIR $WORKDIR
