module V1
  module MediaTypes
    extend Contexts::DSL

    reader = DefinitionReader.new

    reader.for_all_definitions do |defn|
      service = create_service(defn['name'].delete(' '))
      shape_map = Seahorse::Model::ShapeMap.new(defn['shapes'])

      defn['resources'].each do |res|
        next unless res['shape']
        media_type = create_media_type(service, res['shape'])
        media_type.identifier 'application/json'
        shape = shape_map.shape({ 'shape' => res['shape'] })
        media_type.attributes do
          Util.shape_to_attributes(self, shape)
        end
      end
    end
  end
end
