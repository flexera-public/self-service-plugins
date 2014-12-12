class DefinitionReader
  def for_all_definitions(&block)
    @definition_file_names ||= Dir.glob('design/api_definitions/*.yaml')
    @definition_file_names.each do |file_name|
      @definitions ||= {}
      @definitions[file_name] ||= YAML.load(File.read(file_name))
      block.call(@definitions[file_name])
    end
  end
end
