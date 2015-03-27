# Use this file to define your response templates and traits.
#
# For example, to define a response template:
#   response_template :custom do |media_type:|
#     status 200
#     media_type media_type
#   end
Praxis::ApiDefinition.define do
  trait :versionable do
    headers do
      key "X-Api-Version", String, values: ['1.0'], required: true
    end
  end

  response_template :created do
    status 201
    description 'The configuration is created'
    media_type V1::MediaTypes::Configuration
  end
end
