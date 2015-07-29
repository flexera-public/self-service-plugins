module V1
  class Change
    include Praxis::Controller

    implements V1::ApiResources::Change

    def show(id:, **other_params)
      route53 = V1::Helpers::Aws.get_route53_client

      begin
        change = route53.get_change(id: id)
        response = Praxis::Responses::Ok.new()
        response.body = JSON.pretty_generate(V1::MediaTypes::Change.dump(change.change_info))
        response.headers['Content-Type'] = V1::MediaTypes::Change.identifier
      rescue Aws::Route53::Errors::NoSuchChange => e
        response = Praxis::Responses::NotFound.new()
        response.body = { error: e.message }
      rescue Aws::Route53::Errors::InvalidInput => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.message }
      end
      response
    end
  end
end
