module V0
  class HealthCheck
    include Praxis::Controller

    implements V0::ApiResources::HealthCheck

    def health_check(**params)
      response = Praxis::Responses::Ok.new
      response.body = "OK"
      response
    end
  end
end
