  module MediaTypes
    class Bucket < Praxis::MediaType

      identifier 'vnd.rightscale.bucket'

      attributes do
        attribute :instance, String
      end

      view :default do
      end
    end
  end
