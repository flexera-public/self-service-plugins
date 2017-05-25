FROM ruby:2.3-onbuild

ENV APIC_URL https://apic.local
ENV APIC_USER admin
ENV APIC_PASSWORD admin

ENV APP_ENV production

EXPOSE 4567

CMD ["ruby", "application.rb"]
