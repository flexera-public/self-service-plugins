FROM ruby:2.1

RUN mkdir -p /usr/src/app

WORKDIR /usr/src/app

RUN git clone https://github.com/rightscale/self-service-plugins.git

WORKDIR /usr/src/app/self-service-plugins

RUN git checkout dockerize_dme

WORKDIR /usr/src/app/self-service-plugins/digital-ocean/do_api

RUN bundle install

EXPOSE 8088

# overwrite this with 'CMD []' in a dependent Dockerfile
CMD ["rackup","-p","8088"]

