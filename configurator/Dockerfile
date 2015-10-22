FROM rightscale/ops_ruby21x_build:latest
MAINTAINER Kannan Manickam <kannan@rightscale.com>

EXPOSE 8080

ENV RACK_ENV rightscale_changeme

WORKDIR /srv/cm-configurator
VOLUME /var
CMD ["bundle", "exec", "rainbows", "-c", "./config/rainbows.rb", "-p", "8080"]

# Order source dirs by increasing probability of change
# to promote use of layer cache during image builds.
ADD config.ru Rakefile /srv/cm-configurator/
ADD config /srv/cm-configurator/config
ADD design /srv/cm-configurator/design
ADD app /srv/cm-configurator/app

# "Bundle package" could tamper with timestamps; always add
# dependencies last, even though they don't change frequently
ADD Gemfile* /srv/cm-configurator/
ADD vendor/cache /srv/cm-configurator/vendor/cache

RUN bundle install --deployment --without="development test"

