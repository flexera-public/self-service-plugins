# Infamous snakecase

class String

  def camel_case
    return self if self !~ /_/ && self =~ /[A-Z]+.*/
    split('_').map{|e| e.capitalize}.join
  end

  # @param [String] string
  # @return [String] Returns the underscored version of the given string.
  def underscore
    self.
      gsub(@@irregular_regex) { |word| '_' + @@irregular_inflections[word] }.
      gsub(/([A-Z0-9]+)([A-Z][a-z])/, '\1_\2').
      scan(/[a-z0-9]+|\d+|[A-Z0-9]+[a-z]*/).
      join('_').downcase
  end

  # Keep track of irregularities
  def self.irregular_inflections(hash)
    @@irregular_inflections ||= {}
    @@irregular_inflections.update(hash)
    @@irregular_regex = Regexp.new(@@irregular_inflections.keys.join('|'))
  end

  irregular_inflections({
    'ARNs'        => 'arns',
    'CNAMEs'      => 'cnames',
    'Ec2'         => 'ec2',
    'ElastiCache' => 'elasticache',
    'ETag'        => 'etag',
    'iSCSI'       => 'iscsi',
  })
end
