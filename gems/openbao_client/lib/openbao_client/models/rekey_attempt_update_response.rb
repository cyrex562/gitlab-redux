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
  class RekeyAttemptUpdateResponse
    attr_accessor :backup

    attr_accessor :complete

    attr_accessor :keys

    attr_accessor :keys_base64

    attr_accessor :n

    attr_accessor :nounce

    attr_accessor :pgp_fingerprints

    attr_accessor :progress

    attr_accessor :required

    attr_accessor :started

    attr_accessor :t

    attr_accessor :verification_nonce

    attr_accessor :verification_required

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'backup' => :'backup',
        :'complete' => :'complete',
        :'keys' => :'keys',
        :'keys_base64' => :'keys_base64',
        :'n' => :'n',
        :'nounce' => :'nounce',
        :'pgp_fingerprints' => :'pgp_fingerprints',
        :'progress' => :'progress',
        :'required' => :'required',
        :'started' => :'started',
        :'t' => :'t',
        :'verification_nonce' => :'verification_nonce',
        :'verification_required' => :'verification_required'
      }
    end

    # Returns all the JSON keys this model knows about
    def self.acceptable_attributes
      attribute_map.values
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'backup' => :'Boolean',
        :'complete' => :'Boolean',
        :'keys' => :'Array<String>',
        :'keys_base64' => :'Array<String>',
        :'n' => :'Integer',
        :'nounce' => :'String',
        :'pgp_fingerprints' => :'Array<String>',
        :'progress' => :'Integer',
        :'required' => :'Integer',
        :'started' => :'String',
        :'t' => :'Integer',
        :'verification_nonce' => :'String',
        :'verification_required' => :'Boolean'
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
        fail ArgumentError, "The input argument (attributes) must be a hash in `OpenbaoClient::RekeyAttemptUpdateResponse` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `OpenbaoClient::RekeyAttemptUpdateResponse`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'backup')
        self.backup = attributes[:'backup']
      end

      if attributes.key?(:'complete')
        self.complete = attributes[:'complete']
      end

      if attributes.key?(:'keys')
        if (value = attributes[:'keys']).is_a?(Array)
          self.keys = value
        end
      end

      if attributes.key?(:'keys_base64')
        if (value = attributes[:'keys_base64']).is_a?(Array)
          self.keys_base64 = value
        end
      end

      if attributes.key?(:'n')
        self.n = attributes[:'n']
      end

      if attributes.key?(:'nounce')
        self.nounce = attributes[:'nounce']
      end

      if attributes.key?(:'pgp_fingerprints')
        if (value = attributes[:'pgp_fingerprints']).is_a?(Array)
          self.pgp_fingerprints = value
        end
      end

      if attributes.key?(:'progress')
        self.progress = attributes[:'progress']
      end

      if attributes.key?(:'required')
        self.required = attributes[:'required']
      end

      if attributes.key?(:'started')
        self.started = attributes[:'started']
      end

      if attributes.key?(:'t')
        self.t = attributes[:'t']
      end

      if attributes.key?(:'verification_nonce')
        self.verification_nonce = attributes[:'verification_nonce']
      end

      if attributes.key?(:'verification_required')
        self.verification_required = attributes[:'verification_required']
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
          backup == o.backup &&
          complete == o.complete &&
          keys == o.keys &&
          keys_base64 == o.keys_base64 &&
          n == o.n &&
          nounce == o.nounce &&
          pgp_fingerprints == o.pgp_fingerprints &&
          progress == o.progress &&
          required == o.required &&
          started == o.started &&
          t == o.t &&
          verification_nonce == o.verification_nonce &&
          verification_required == o.verification_required
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [backup, complete, keys, keys_base64, n, nounce, pgp_fingerprints, progress, required, started, t, verification_nonce, verification_required].hash
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
