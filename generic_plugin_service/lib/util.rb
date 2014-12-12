class Util
  # Converts a shape definition to a attributor type and options
  #
  def self.shape_to_attributes(context, shape)
    shape.members.each do |name, detail|
      options = discover_options(detail)
      case detail.type
      ## Basic Types
      when 'string', 'character', 'byte'
        context.attribute name.to_sym, Attributor::String, options
      when 'integer', 'long'
        context.attribute name.to_sym, Attributor::Integer, options
      when 'float', 'double'
        context.attribute name.to_sym Attributor::Float, options
      when 'boolean'
        context.attribute name.to_sym, Attributor::Boolean, options
      ## Complex Types
      when 'structure'
        # TODO: Do it completely
        context.attribute name.to_sym, Attributor::Struct do
          Util.shape_to_attributes(self, detail)
        end
      when 'list'
        if detail.member.type == 'structure'
          # TODO: Fix it
          next
          context.attribute name.to_sym, Attributor::Collection.of(Attributor::Struct) do
            Util.shape_to_attributes(self, detail.member)
          end
        else
          context.attribute name.to_sym, Attributor::Collection.of(discover_type(detail.member))
        end
      when 'map'
        key_type = discover_type(detail.key)
        value_type = discover_type(detail.value)
        context.attribute name.to_sym, Attributor::Hash.of(key: key_type, value: value_type)
      when 'timestamp'
        # TODO: Do we need to handle anything special here?
        context.attribute name.to_sym, Attributor::DateTime, options
      end
    end
  end

  def self.discover_type(shape)
    case shape.type
    when 'string', 'character', 'byte'
      Attributor::String
    when 'integer', 'long'
      Attributor::Integer
    when 'float', 'double'
      Attributor::Float
    when 'boolean'
      Attributor::Boolean
    end
  end

  def self.discover_options(detail)
    options = {}
    case detail.type
    when 'string'
      options[:min] = detail.min if detail.min
      options[:max] = detail.max if detail.max
      options[:values] = detail.enum.to_a if detail.enum
      options[:regexp] = detail.pattern if detail.pattern
    end
    options
  end
end
