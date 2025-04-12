# frozen_string_literal: true

module Packages
  module Maven
    class CreatePackageService < ::Packages::CreatePackageService
      def execute
        return ERROR_RESPONSE_PACKAGE_PROTECTED if package_protected?

        app_group, _, app_name = params[:name].rpartition('/')
        app_group.tr!('/', '.')

        package = create_package!(:maven,
          maven_metadatum_attributes: {
            path: params[:path],
            app_group: app_group,
            app_name: app_name,
            app_version: params[:version]
          }
        )

        ServiceResponse.success(payload: { package: package })
      rescue ActiveRecord::RecordInvalid => e
        reason = e.record&.errors&.of_kind?(:name, :taken) ? :name_taken : :invalid_parameter

        ServiceResponse.error(message: e.message, reason: reason)
      end

      private

      def package_protected?
        return false if Feature.disabled?(:packages_protected_packages_maven, project)

        super(package_name: params[:name], package_type: :maven)
      end
    end
  end
end
