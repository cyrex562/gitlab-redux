# frozen_string_literal: true

module Types
  module PermissionTypes
    module Namespaces
      class Base < BasePermissionType
        graphql_name 'NamespacePermissions'

        ability_field :read_namespace
      end
    end
  end
end

::Types::PermissionTypes::Namespaces::Base.prepend_mod
