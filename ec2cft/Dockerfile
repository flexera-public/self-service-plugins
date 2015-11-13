FROM  ruby:2.1-onbuild 

RUN git clone https://github.com/rightscale/self-service-plugins.git

WORKDIR /usr/src/app/self-service-plugins/ec2cft

RUN git checkout dockerize_dme

RUN bundle install

EXPOSE 8088

# overwrite this with 'CMD []' in a dependent Dockerfile
CMD ["bundle","exec","rackup","-p","8088","-o","0.0.0.0"]

