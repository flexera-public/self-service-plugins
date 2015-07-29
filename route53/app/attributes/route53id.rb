module Attributes
  class Route53Id
    include Attributor::Type

    def self.native_type
      return ::String
    end

    def self.load(value,context=Attributor::DEFAULT_ROOT_CONTEXT, **options)
      if value.kind_of?(Enumerable)
        raise IncompatibleTypeError,  context: context, value_type: value.class, type: self
      end

      # This is my hack to only get the ID from the path, this whole thing is just
      # copy/pasted from https://github.com/rightscale/attributor/blob/master/lib/attributor/types/string.rb
      # might be completely wrong and awful
      value.match(/\/[a-z_]*\/([a-z0-9A-Z_]*)$/)[1] && String(value).match(/\/[a-z_]*\/([a-z0-9A-Z_]*)$/)[1]
    rescue
      super
    end

    def self.example(context=nil, options:{})
      if options[:regexp]
        # It may fail to generate an example, see bug #72.
        options[:regexp].gen rescue ('Failed to generate example for %s' % options[:regexp].inspect)
      else
        /\w+/.gen
      end
    end

    def self.family
      'string'
    end
  end
end
