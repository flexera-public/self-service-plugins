FROM ruby:2.2
MAINTAINER Ryan J. Geyer <me@ryangeyer.com>

RUN gem install praxis -v 0.16.1

ENV WORKDIR /opt/praxis
COPY app $WORKDIR/app
COPY config $WORKDIR/config
COPY design $WORKDIR/design
COPY config.ru $WORKDIR
COPY Gemfile $WORKDIR
COPY Gemfile.lock $WORKDIR
COPY Rakefile $WORKDIR
WORKDIR $WORKDIR

RUN cd $WORKDIR && bundle install
RUN apt-get update && apt-get install -y vim
