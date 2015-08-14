module V1
  class Record
    include Praxis::Controller

    implements V1::ApiResources::Record

    def get_resource_record_set_request(hosted_zone_id, action, name, type, ttl, values)
      resource_set_request = {
        hosted_zone_id: hosted_zone_id,
        change_batch: {
          changes: [
            {
              action: action,
              resource_record_set: {
                name: name,
                type: type,
                ttl: ttl,
                resource_records: []
              }
            }
          ]
        }
      }

      record_values = []

      values.each do |value|
        record_values << { value: value }
      end

      resource_set_request[:change_batch][:changes].first[:resource_record_set][:resource_records] = record_values
      resource_set_request
    end

    def get_resource_record_sets_as_record_model(route53, public_zone_id)
      records_response = route53.list_resource_record_sets(hosted_zone_id: public_zone_id)
      records = []
      records_response.resource_record_sets.each do |record_set|
        records << V1::Models::Record.new(public_zone_id, record_set)
      end
      records
    end

    def do_delete(public_zone_id, id, return_change=false)
      route53 = V1::Helpers::Aws.get_route53_client

      response = return_change ? Praxis::Responses::Ok.new : Praxis::Responses::NoContent.new()
      all_records = nil

      href = V1::ApiResources::PublicZone.prefix+'/'+public_zone_id+
        V1::ApiResources::Record.prefix+'/'+id

      begin
        all_records = get_resource_record_sets_as_record_model(route53, public_zone_id)
      rescue Aws::Route53::Errors::InvalidInput => e
        response = Praxis::Responses::BadRequest.new()
        response.body = {
          error: "While fetching all zone records to locate #{href}\n#{e.inspect}"
        }
      rescue Aws::Route53::Errors::NoSuchHostedZone => e
        response = Praxis::Responses::NotFound.new()
        response.body = {
          error: "While fetching all zone records to locate #{href}\n#{e.inspect}"
        }
      end

      if all_records
        records = all_records.select{|r| r.id == id }
        if records.size > 0
          begin
            record = records.first
            resource_set_request = get_resource_record_set_request(
              public_zone_id,
              "DELETE",
              record.name,
              record.type,
              record.ttl,
              record.values
            )

            aws_response = route53.change_resource_record_sets(resource_set_request)
            if return_change
              response.body = JSON.pretty_generate(V1::MediaTypes::Change.render(aws_response.change_info))
              response.headers['Content-Type'] = V1::MediaTypes::Change.identifier
            end
          rescue  Aws::Route53::Errors::NoSuchHostedZone,
                  Aws::Route53::Errors::NoSuchHealthCheck => e
            response = Praxis::Responses::NotFound.new()
            response.body = { error: e.inspect }
          rescue  Aws::Route53::Errors::PriorRequestNotComplete,
                  Aws::Route53::Errors::InvalidInput,
                  Aws::Route53::Errors::InvalidChangeBatch => e
            response = Praxis::Responses::BadRequest.new()
            response.body = { error: e.inspect }
          end
        else
          response = Praxis::Responses::NotFound.new()
          response.body = {
            error: "Could not find record #{href}"
          }
        end
      end
      response
    end

    def index(public_zone_id:, **params)
      route53 = V1::Helpers::Aws.get_route53_client

      response = self.response

      begin
        records = get_resource_record_sets_as_record_model(route53, public_zone_id)
        records_mediatype = records.map{|r| V1::MediaTypes::Record.render(r) }

        response.body = JSON.pretty_generate(records_mediatype)
        response.headers['Content-Type'] = V1::MediaTypes::Record.identifier+';type=collection'
      rescue Aws::Route53::Errors::InvalidInput => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.inspect }
      rescue Aws::Route53::Errors::NoSuchHostedZone => e
        response = Praxis::Responses::NotFound.new()
        response.body = { error: e.inspect }
      end

      response
    end

    def show(public_zone_id:, id:, **other_params)
      route53 = V1::Helpers::Aws.get_route53_client

      response = self.response

      begin
        records = get_resource_record_sets_as_record_model(route53, public_zone_id)
        records_mediatype = records.map{|r| V1::MediaTypes::Record.render(r) }


        record = records_mediatype.select{|r| r[:id] == id }

        if record.size == 0
          response = Praxis::Responses::NotFound.new()
          response.body = { error: "Record ID (#{id}) not found in zone ID (#{public_zone_id})"}
        else
          response.body = JSON.pretty_generate(record.first)
          response.headers['Content-Type'] = V1::MediaTypes::Record.identifier
        end
      rescue Aws::Route53::Errors::NoSuchHostedZone => e
        response = Praxis::Responses::NotFound.new()
        response.body = { error: e.inspect }
      rescue Aws::Route53::Errors::InvalidInput => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.inspect }
      end
      response
    end

    def create(**other_params)
      route53 = V1::Helpers::Aws.get_route53_client

      response = Praxis::Responses::Created.new()

      public_zone_id = ''

      if request.params && request.params.public_zone_id
        public_zone_id = request.params.public_zone_id
      else
        public_zone_id = request.payload.public_zone_id
      end

      begin
        resource_set_request = get_resource_record_set_request(
          public_zone_id,
          "UPSERT",
          request.payload.name,
          request.payload.type,
          request.payload.ttl,
          request.payload.values
        )

        aws_response = route53.change_resource_record_sets(resource_set_request)
        name = request.payload.name
        if name[-1, 1] != '.'
          name = request.payload.name+'.'
        end
        record_hash = {
          name: name,
          type: request.payload.type,
          ttl: request.payload.ttl,
          resource_records: request.payload.values.map{|r| OpenStruct.new({ value: r }) }
        }
        record_struct = OpenStruct.new(record_hash)
        record = V1::Models::Record.new(public_zone_id, record_struct)
        response.body = JSON.pretty_generate(V1::MediaTypes::Change.render(aws_response.change_info))
        response.headers['Content-Type'] = V1::MediaTypes::Change.identifier
        response.headers['Location'] = record.href
      rescue  Aws::Route53::Errors::NoSuchHostedZone,
              Aws::Route53::Errors::NoSuchHealthCheck => e
        response = Praxis::Responses::NotFound.new()
        response.body = { error: e.inspect }
      rescue  Aws::Route53::Errors::PriorRequestNotComplete,
              Aws::Route53::Errors::InvalidInput,
              Aws::Route53::Errors::InvalidChangeBatch => e
        response = Praxis::Responses::BadRequest.new()
        response.body = { error: e.inspect }
      end

      response
    end

    def delete(public_zone_id:, id:, **other_params)
      response = do_delete(public_zone_id, id)
    end

    def release(public_zone_id:, id:, **other_params)
      response = do_delete(public_zone_id, id, true)
    end

  end
end
