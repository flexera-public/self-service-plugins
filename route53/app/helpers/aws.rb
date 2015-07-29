module V1
  module Helpers
    module Aws

      def self.get_route53_client()
        credentials = ::Aws::Credentials.new(ENV['AWS_ACCESS_KEY_ID'], ENV['AWS_SECRET_ACCESS_KEY'])
        route53 = ::Aws::Route53::Client.new(region: 'us-east-1', credentials: credentials)
      end

    end
  end
end
