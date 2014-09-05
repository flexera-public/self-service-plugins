  module MediaTypes
    class Instance < Praxis::MediaType

      identifier 'vnd.rightscale.instance' #+json;type=collection'

      attributes do
        attribute :id, String
        attribute :state, String
        attribute :region, String
        attribute :databaseVersion, String
        attribute :currentDiskSize, Integer
        attribute :maxDiskSize, Integer
        attribute :ipAddresses, String #Attributor::Collection.of(Struct) do
#          attribute :ipAddress, String
#          attribute :timeToRetire, DateTime
#        end

        attribute :instanceType, String
        attribute :masterInstanceName, String
        attribute :replicaNames, String

      end

      view :default do
        attribute :id
        attribute :state
        attribute :region
        attribute :databaseVersion
        attribute :currentDiskSize
        attribute :maxDiskSize
        attribute :ipAddresses
        attribute :instanceType
        attribute :masterInstanceName
        attribute :replicaNames
      end
    end
  end
