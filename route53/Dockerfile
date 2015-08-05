FROM ruby:2.2
MAINTAINER Ryan J. Geyer <me@ryangeyer.com>

RUN gem install praxis -v 0.16.1

COPY Gemfile /tmp/Gemfile
RUN cd /tmp && bundle install

ENV WORKDIR /opt/praxis
WORKDIR $WORKDIR
