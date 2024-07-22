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
  class JwtConfigureRequest
    # The value against which to match the 'iss' claim in a JWT. Optional.
    attr_accessor :bound_issuer

    # The default role to use if none is provided during login. If not set, a role is required during login.
    attr_accessor :default_role

    # The CA certificate or chain of certificates, in PEM format, to use to validate connections to the JWKS URL. If not set, system certificates are used.
    attr_accessor :jwks_ca_pem

    # JWKS URL to use to authenticate signatures. Cannot be used with \"oidc_discovery_url\" or \"jwt_validation_pubkeys\".
    attr_accessor :jwks_url

    # A list of supported signing algorithms. Defaults to RS256.
    attr_accessor :jwt_supported_algs

    # A list of PEM-encoded public keys to use to authenticate signatures locally. Cannot be used with \"jwks_url\" or \"oidc_discovery_url\".
    attr_accessor :jwt_validation_pubkeys

    # Pass namespace in the OIDC state parameter instead of as a separate query parameter. With this setting, the allowed redirect URL(s) in OpenBao and on the provider side should not contain a namespace query parameter. This means only one redirect URL entry needs to be maintained on the provider side for all OpenBao namespaces that will be authenticating against it. Defaults to true for new configs.
    attr_accessor :namespace_in_state

    # The OAuth Client ID configured with your OIDC provider.
    attr_accessor :oidc_client_id

    # The OAuth Client Secret configured with your OIDC provider.
    attr_accessor :oidc_client_secret

    # The CA certificate or chain of certificates, in PEM format, to use to validate connections to the OIDC Discovery URL. If not set, system certificates are used.
    attr_accessor :oidc_discovery_ca_pem

    # OIDC Discovery URL, without any .well-known component (base path). Cannot be used with \"jwks_url\" or \"jwt_validation_pubkeys\".
    attr_accessor :oidc_discovery_url

    # The response mode to be used in the OAuth2 request. Allowed values are 'query' and 'form_post'.
    attr_accessor :oidc_response_mode

    # The response types to request. Allowed values are 'code' and 'id_token'. Defaults to 'code'.
    attr_accessor :oidc_response_types

    # Provider-specific configuration. Optional.
    attr_accessor :provider_config

    # Attribute mapping from ruby-style variable name to JSON key.
    def self.attribute_map
      {
        :'bound_issuer' => :'bound_issuer',
        :'default_role' => :'default_role',
        :'jwks_ca_pem' => :'jwks_ca_pem',
        :'jwks_url' => :'jwks_url',
        :'jwt_supported_algs' => :'jwt_supported_algs',
        :'jwt_validation_pubkeys' => :'jwt_validation_pubkeys',
        :'namespace_in_state' => :'namespace_in_state',
        :'oidc_client_id' => :'oidc_client_id',
        :'oidc_client_secret' => :'oidc_client_secret',
        :'oidc_discovery_ca_pem' => :'oidc_discovery_ca_pem',
        :'oidc_discovery_url' => :'oidc_discovery_url',
        :'oidc_response_mode' => :'oidc_response_mode',
        :'oidc_response_types' => :'oidc_response_types',
        :'provider_config' => :'provider_config'
      }
    end

    # Returns all the JSON keys this model knows about
    def self.acceptable_attributes
      attribute_map.values
    end

    # Attribute type mapping.
    def self.openapi_types
      {
        :'bound_issuer' => :'String',
        :'default_role' => :'String',
        :'jwks_ca_pem' => :'String',
        :'jwks_url' => :'String',
        :'jwt_supported_algs' => :'Array<String>',
        :'jwt_validation_pubkeys' => :'Array<String>',
        :'namespace_in_state' => :'Boolean',
        :'oidc_client_id' => :'String',
        :'oidc_client_secret' => :'String',
        :'oidc_discovery_ca_pem' => :'String',
        :'oidc_discovery_url' => :'String',
        :'oidc_response_mode' => :'String',
        :'oidc_response_types' => :'Array<String>',
        :'provider_config' => :'Object'
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
        fail ArgumentError, "The input argument (attributes) must be a hash in `OpenbaoClient::JwtConfigureRequest` initialize method"
      end

      # check to see if the attribute exists and convert string to symbol for hash key
      attributes = attributes.each_with_object({}) { |(k, v), h|
        if (!self.class.attribute_map.key?(k.to_sym))
          fail ArgumentError, "`#{k}` is not a valid attribute in `OpenbaoClient::JwtConfigureRequest`. Please check the name to make sure it's valid. List of attributes: " + self.class.attribute_map.keys.inspect
        end
        h[k.to_sym] = v
      }

      if attributes.key?(:'bound_issuer')
        self.bound_issuer = attributes[:'bound_issuer']
      end

      if attributes.key?(:'default_role')
        self.default_role = attributes[:'default_role']
      end

      if attributes.key?(:'jwks_ca_pem')
        self.jwks_ca_pem = attributes[:'jwks_ca_pem']
      end

      if attributes.key?(:'jwks_url')
        self.jwks_url = attributes[:'jwks_url']
      end

      if attributes.key?(:'jwt_supported_algs')
        if (value = attributes[:'jwt_supported_algs']).is_a?(Array)
          self.jwt_supported_algs = value
        end
      end

      if attributes.key?(:'jwt_validation_pubkeys')
        if (value = attributes[:'jwt_validation_pubkeys']).is_a?(Array)
          self.jwt_validation_pubkeys = value
        end
      end

      if attributes.key?(:'namespace_in_state')
        self.namespace_in_state = attributes[:'namespace_in_state']
      end

      if attributes.key?(:'oidc_client_id')
        self.oidc_client_id = attributes[:'oidc_client_id']
      end

      if attributes.key?(:'oidc_client_secret')
        self.oidc_client_secret = attributes[:'oidc_client_secret']
      end

      if attributes.key?(:'oidc_discovery_ca_pem')
        self.oidc_discovery_ca_pem = attributes[:'oidc_discovery_ca_pem']
      end

      if attributes.key?(:'oidc_discovery_url')
        self.oidc_discovery_url = attributes[:'oidc_discovery_url']
      end

      if attributes.key?(:'oidc_response_mode')
        self.oidc_response_mode = attributes[:'oidc_response_mode']
      end

      if attributes.key?(:'oidc_response_types')
        if (value = attributes[:'oidc_response_types']).is_a?(Array)
          self.oidc_response_types = value
        end
      end

      if attributes.key?(:'provider_config')
        self.provider_config = attributes[:'provider_config']
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
          bound_issuer == o.bound_issuer &&
          default_role == o.default_role &&
          jwks_ca_pem == o.jwks_ca_pem &&
          jwks_url == o.jwks_url &&
          jwt_supported_algs == o.jwt_supported_algs &&
          jwt_validation_pubkeys == o.jwt_validation_pubkeys &&
          namespace_in_state == o.namespace_in_state &&
          oidc_client_id == o.oidc_client_id &&
          oidc_client_secret == o.oidc_client_secret &&
          oidc_discovery_ca_pem == o.oidc_discovery_ca_pem &&
          oidc_discovery_url == o.oidc_discovery_url &&
          oidc_response_mode == o.oidc_response_mode &&
          oidc_response_types == o.oidc_response_types &&
          provider_config == o.provider_config
    end

    # @see the `==` method
    # @param [Object] Object to be compared
    def eql?(o)
      self == o
    end

    # Calculates hash code according to all attributes.
    # @return [Integer] Hash code
    def hash
      [bound_issuer, default_role, jwks_ca_pem, jwks_url, jwt_supported_algs, jwt_validation_pubkeys, namespace_in_state, oidc_client_id, oidc_client_secret, oidc_discovery_ca_pem, oidc_discovery_url, oidc_response_mode, oidc_response_types, provider_config].hash
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
