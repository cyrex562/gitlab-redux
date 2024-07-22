=begin
#OpenBao API

#HTTP API that gives you full access to OpenBao. All API routes are prefixed with `/v1/`.

The version of the OpenAPI document: 2.0.0

Generated by: https://openapi-generator.tech
Generator version: 7.7.0

=end

require 'date'
require 'time'

module OpenbaoClient
  class TokenCreateRequest
    # Name to associate with this token
    attr_accessor :display_name

    # Name of the entity alias to associate with this token
    attr_accessor :entity_alias

    # Explicit Max TTL of this token
    attr_accessor :explicit_max_ttl

    # Value for the token
    attr_accessor :id

    # Use 'ttl' instead
    attr_accessor :lease

    # Arbitrary key=value metadata to associate with the token
    attr_accessor :meta

    # Do not include default policy for this token
    attr_accessor :no_default_policy

    # Create the token with no parent
    attr_accessor :no_parent

    # Max number of uses for this token
    attr_accessor :num_uses

    # Renew period
    attr_accessor :period

    # List of policies for the token
    attr_accessor :policies

    # Allow token to be renewed past its initial TTL up to system/mount maximum TTL
    attr_accessor :renewable

    # Time to live for this token
    attr_accessor :ttl

    # Token type
    attr_accessor :type

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'display_name' => :'display_name',
        :'entity_alias' => :'entity_alias',
        :'explicit_max_ttl' => :'explicit_max_ttl',
        :'id' => :'id',
        :'lease' => :'lease',
        :'meta' => :'meta',
        :'no_default_policy' => :'no_default_policy',
        :'no_parent' => :'no_parent',
        :'num_uses' => :'num_uses',
        :'period' => :'period',
        :'policies' => :'policies',
        :'renewable' => :'renewable',
        :'ttl' => :'ttl',
        :'type' => :'type'
      }
    end

    # Returns all the JSON keys this model knows about
    def self.acceptable_attributes
      attribute_map.values
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'display_name' => :'String',
        :'entity_alias' => :'String',
        :'explicit_max_ttl' => :'String',
        :'id' => :'String',
        :'lease' => :'String',
        :'meta' => :'Object',
        :'no_default_policy' => :'Boolean',
        :'no_parent' => :'Boolean',
        :'num_uses' => :'Integer',
        :'period' => :'String',
        :'policies' => :'Array<String>',
        :'renewable' => :'Boolean',
        :'ttl' => :'String',
        :'type' => :'String'
      }
    end

    # List of attributes with nullable: true
    def self.openapi_nullable
      Set.new([
      ])
    end

    # Initializes the object
    # @param [Hash] attributes Model attributes in the form of hash
    def initialize(attributes = {})
      if (!attributes.is_a?(Hash))
        fail ArgumentError, "The input argument (attributes) must be a hash in `OpenbaoClient::TokenCreateRequest` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `OpenbaoClient::TokenCreateRequest`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'display_name')
        self.display_name = attributes[:'display_name']
      end

      if attributes.key?(:'entity_alias')
        self.entity_alias = attributes[:'entity_alias']
      end

      if attributes.key?(:'explicit_max_ttl')
        self.explicit_max_ttl = attributes[:'explicit_max_ttl']
      end

      if attributes.key?(:'id')
        self.id = attributes[:'id']
      end

      if attributes.key?(:'lease')
        self.lease = attributes[:'lease']
      end

      if attributes.key?(:'meta')
        self.meta = attributes[:'meta']
      end

      if attributes.key?(:'no_default_policy')
        self.no_default_policy = attributes[:'no_default_policy']
      end

      if attributes.key?(:'no_parent')
        self.no_parent = attributes[:'no_parent']
      end

      if attributes.key?(:'num_uses')
        self.num_uses = attributes[:'num_uses']
      end

      if attributes.key?(:'period')
        self.period = attributes[:'period']
      end

      if attributes.key?(:'policies')
        if (value = attributes[:'policies']).is_a?(Array)
          self.policies = value
        end
      end

      if attributes.key?(:'renewable')
        self.renewable = attributes[:'renewable']
      else
        self.renewable = true
      end

      if attributes.key?(:'ttl')
        self.ttl = attributes[:'ttl']
      end

      if attributes.key?(:'type')
        self.type = attributes[:'type']
      end
    end

    # Show invalid properties with the reasons. Usually used together with valid?
    # @return Array for valid properties with the reasons
    def list_invalid_properties
      warn '[DEPRECATED] the `list_invalid_properties` method is obsolete'
      invalid_properties = Array.new
      invalid_properties
    end

    # Check to see if the all the properties in the model are valid
    # @return true if the model is valid
    def valid?
      warn '[DEPRECATED] the `valid?` method is obsolete'
      true
    end

    # Checks equality by comparing each attribute.
    # @param [Object] Object to be compared
    def ==(o)
      return true if self.equal?(o)
      self.class == o.class &&
          display_name == o.display_name &&
          entity_alias == o.entity_alias &&
          explicit_max_ttl == o.explicit_max_ttl &&
          id == o.id &&
          lease == o.lease &&
          meta == o.meta &&
          no_default_policy == o.no_default_policy &&
          no_parent == o.no_parent &&
          num_uses == o.num_uses &&
          period == o.period &&
          policies == o.policies &&
          renewable == o.renewable &&
          ttl == o.ttl &&
          type == o.type
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [display_name, entity_alias, explicit_max_ttl, id, lease, meta, no_default_policy, no_parent, num_uses, period, policies, renewable, ttl, type].hash
    end

    # Builds the object from hash
    # @param [Hash] attributes Model attributes in the form of hash
    # @return [Object] Returns the model itself
    def self.build_from_hash(attributes)
      return nil unless attributes.is_a?(Hash)
      attributes = attributes.transform_keys(&:to_sym)
      transformed_hash = {}
      openapi_types.each_pair do |key, type|
        if attributes.key?(attribute_map[key]) && attributes[attribute_map[key]].nil?
          transformed_hash["#{key}"] = nil
        elsif type =~ /\AArray<(.*)>/i
          # check to ensure the input is an array given that the attribute
          # is documented as an array but the input is not
          if attributes[attribute_map[key]].is_a?(Array)
            transformed_hash["#{key}"] = attributes[attribute_map[key]].map { |v| _deserialize($1, v) }
          end
        elsif !attributes[attribute_map[key]].nil?
          transformed_hash["#{key}"] = _deserialize(type, attributes[attribute_map[key]])
        end
      end
      new(transformed_hash)
    end

    # Deserializes the data based on type
    # @param string type Data type
    # @param string value Value to be deserialized
    # @return [Object] Deserialized data
    def self._deserialize(type, value)
      case type.to_sym
      when :Time
        Time.parse(value)
      when :Date
        Date.parse(value)
      when :String
        value.to_s
      when :Integer
        value.to_i
      when :Float
        value.to_f
      when :Boolean
        if value.to_s =~ /\A(true|t|yes|y|1)\z/i
          true
        else
          false
        end
      when :Object
        # generic object (usually a Hash), return directly
        value
      when /\AArray<(?<inner_type>.+)>\z/
        inner_type = Regexp.last_match[:inner_type]
        value.map { |v| _deserialize(inner_type, v) }
      when /\AHash<(?<k_type>.+?), (?<v_type>.+)>\z/
        k_type = Regexp.last_match[:k_type]
        v_type = Regexp.last_match[:v_type]
        {}.tap do |hash|
          value.each do |k, v|
            hash[_deserialize(k_type, k)] = _deserialize(v_type, v)
          end
        end
      else # model
        # models (e.g. Pet) or oneOf
        klass = OpenbaoClient.const_get(type)
        klass.respond_to?(:openapi_any_of) || klass.respond_to?(:openapi_one_of) ? klass.build(value) : klass.build_from_hash(value)
      end
    end

    # Returns the string representation of the object
    # @return [String] String presentation of the object
    def to_s
      to_hash.to_s
    end

    # to_body is an alias to to_hash (backward compatibility)
    # @return [Hash] Returns the object in the form of hash
    def to_body
      to_hash
    end

    # Returns the object in the form of hash
    # @return [Hash] Returns the object in the form of hash
    def to_hash
      hash = {}
      self.class.attribute_map.each_pair do |attr, param|
        value = self.send(attr)
        if value.nil?
          is_nullable = self.class.openapi_nullable.include?(attr)
          next if !is_nullable || (is_nullable && !instance_variable_defined?(:"@#{attr}"))
        end

        hash[param] = _to_hash(value)
      end
      hash
    end

    # Outputs non-array value in the form of hash
    # For object, use to_hash. Otherwise, just return the value
    # @param [Object] value Any valid value
    # @return [Hash] Returns the value in the form of hash
    def _to_hash(value)
      if value.is_a?(Array)
        value.compact.map { |v| _to_hash(v) }
      elsif value.is_a?(Hash)
        {}.tap do |hash|
          value.each { |k, v| hash[k] = _to_hash(v) }
        end
      elsif value.respond_to? :to_hash
        value.to_hash
      else
        value
      end
    end

  end

end
