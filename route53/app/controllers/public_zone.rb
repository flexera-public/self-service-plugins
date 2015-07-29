require 'securerandom'

module V1
  class PublicZone
    include Praxis::Controller

    implements V1::ApiResources::PublicZone

    def index(**params)
      route53 = V1::Helpers::Aws.get_route53_client

      # Deliberately not catching SelegationSet errors since it's not currently
      # supported in this plugin
      # https://github.com/aws/aws-sdk-ruby/blob/9255277a1da95a6217f603e683bd49cc677a4b5a/aws-sdk-core/apis/route53/2013-04-01/api-2.json#L612-L626
      begin
        zones_mediatype = []
        list_hosted_zones_response = route53.list_hosted_zones

        list_hosted_zones_response.hosted_zones.each do |native_zone|
          zone = V1::Models::PublicZone.new(native_zone)
          zones_mediatype << V1::MediaTypes::PublicZone.render(zone)
        end

        response = Praxis::Responses::Ok.new()
        response.body = JSON.pretty_generate(zones_mediatype)
        response.headers['Content-Type'] = V1::MediaTypes::PublicZone.identifier+';type=collection'
      rescue Aws::Route53::Errors::InvalidInput => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.message }
      end

      response
    end

    def show(id:, **other_params)
      route53 = V1::Helpers::Aws.get_route53_client

      # https://github.com/aws/aws-sdk-ruby/blob/9255277a1da95a6217f603e683bd49cc677a4b5a/aws-sdk-core/apis/route53/2013-04-01/api-2.json#L514-L525
      begin
        zone = route53.get_hosted_zone(id: id)
        response = Praxis::Responses::Ok.new()
        zone_hash = V1::Models::PublicZone.new(zone.hosted_zone)
        response.body = JSON.pretty_generate(V1::MediaTypes::PublicZone.render(zone_hash))
        response.headers['Content-Type'] = V1::MediaTypes::PublicZone.identifier
      rescue Aws::Route53::Errors::NoSuchHostedZone => e
        response = Praxis::Responses::NotFound.new()
        response.body = { error: e.message }
      rescue Aws::Route53::Errors::InvalidInput => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.message }
      end
      response
    end

    def create(**other_params)
      route53 = V1::Helpers::Aws.get_route53_client

      zone_params = {
        name: request.payload.name,
        caller_reference: SecureRandom.uuid
      }

      begin
        aws_response = route53.create_hosted_zone(zone_params)

        response = Praxis::Responses::Created.new()
        zone_model = V1::Models::PublicZone.new(aws_response.hosted_zone, aws_response.change_info)
        # zone_shaped_hash = Hash[aws_response.hosted_zone]
        # zone_shaped_hash[:change] = aws_response.change_info
        # zone = V1::MediaTypes::PublicZone.render(zone_shaped_hash)
        zone = V1::MediaTypes::PublicZone.render(zone_model)
        response.body = JSON.pretty_generate(zone)
        response.headers['Content-Type'] = V1::MediaTypes::PublicZone.identifier
        response.headers['Location'] = zone[:href]
      rescue  Aws::Route53::Errors::ConflictingDomainExists,
              Aws::Route53::Errors::InvalidInput,
              Aws::Route53::Errors::TooManyHostedZones,
              Aws::Route53::Errors::HostedZoneAlreadyExists => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.message }
      end

      response
    end

    def delete(id:, **other_params)
      route53 = V1::Helpers::Aws.get_route53_client

      begin
        delete_response = route53.delete_hosted_zone(id: id)
        response = Praxis::Responses::Ok.new()
        response.body = JSON.pretty_generate(V1::MediaTypes::Change.render(delete_response.change_info))
        response.headers['Content-Type'] = V1::MediaTypes::Change.identifier
      rescue Aws::Route53::Errors::NoSuchHostedZone => e
        response = Praxis::Responses::NotFound.new()
        response.body = { error: e.message }
      rescue  Aws::Route53::Errors::InvalidInput,
              Aws::Route53::Errors::PriorRequestNotComplete,
              Aws::Route53::Errors::HostedZoneNotEmpty => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.message }
      end

      response
    end

  end
end
