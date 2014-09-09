  module MediaTypes
    class Instance < Praxis::MediaType

      identifier 'vnd.rightscale.instance'

      attributes do
        attribute :instance, String
        attribute :state, String
        attribute :region, String
        attribute :databaseVersion, String
        attribute :currentDiskSize, Integer
        attribute :maxDiskSize, Integer
        attribute :ipAddresses, Attributor::Collection.of(Attributor::Struct) do
          attribute :ipAddress, String
          attribute :timeToRetire, DateTime
        end

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
