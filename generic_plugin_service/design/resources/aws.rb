module V1
  module ApiResources
    extend Contexts::DSL
    # Mapping of actions and verbs for some basic actions. AWS uses POST
    # for everything!
    VERB_MAPPINGS = {
      'create' => :post,
      'index'  => :get,
      'show'   => :get,
      'update' => :put,
      'delete' => :delete
    }

    reader = DefinitionReader.new
    reader.for_all_definitions do |defn|
      # Create a module with the name of the service
      service = create_service(defn['name'].delete(' '))
      defn['resources'].each do |res|
        # Create a resource
        resource = create_resource(service, res['name'])

        resource.media_type V1::MediaTypes.const_get(
          defn['name'].delete(' ').camel_case).const_get(res['shape'].camel_case)

        resource.version defn['version']
        resource.routing do
          prefix "/api#{defn['url']}/#{res['name'].pluralize}"
        end

        res['actions'] && res['actions'].each do |act|
          action_name = act['name']
          resource.action(action_name) do
            use :versionable

            routing do
              send(VERB_MAPPINGS[action_name] || :post, act['path'])
            end

            shape_map = Seahorse::Model::ShapeMap.new(defn['shapes'])
            if act['payload']
              shape = shape_map.shape({ 'shape' => act['payload'] })
              payload do
                Util.shape_to_attributes(self, shape)
              end
            end

            response :ok
          end
        end
      end
    end
  end
end
