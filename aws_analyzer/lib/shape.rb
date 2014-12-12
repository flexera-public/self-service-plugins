require 'seahorse'

module AwsAnalyzer

  # Description of request or response body structures
  class Shape < Seahorse::Model::Shapes::Shape

    # Hash representation
    def to_hash
      @definition
    end

    # YAML representation
    def to_yaml
      YAML.dump(to_hash)
    end
    alias :to_s :to_yaml

  end

end

