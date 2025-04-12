# frozen_string_literal: true

require 'spec_helper'

RSpec.describe API::Conan::V2::ProjectPackages, feature_category: :package_registry do
  include_context 'with conan api setup'

  let_it_be_with_reload(:package) { create(:conan_package, project: project) }
  let(:project_id) { project.id }
  let(:url) { "/projects/#{project_id}/packages/conan/v2/conans/#{url_suffix}" }

  shared_examples 'conan package revisions feature flag check' do
    before do
      stub_feature_flags(conan_package_revisions_support: false)
    end

    it_behaves_like 'returning response status with message', status: :not_found,
      message: "404 'conan_package_revisions_support' feature flag is disabled Not Found"
  end

  shared_examples 'packages feature check' do
    before do
      stub_packages_setting(enabled: false)
    end

    it_behaves_like 'returning response status', :not_found
  end

  describe 'GET /api/v4/projects/:id/packages/conan/v2/users/check_credentials' do
    let(:url) { "/projects/#{project.id}/packages/conan/v2/users/check_credentials" }

    it_behaves_like 'conan check_credentials endpoint'
  end

  describe 'GET /api/v4/projects/:id/packages/conan/v2/conans/search' do
    let(:url) { "/projects/#{project.id}/packages/conan/v2/conans/search" }

    it_behaves_like 'conan search endpoint'

    it_behaves_like 'conan FIPS mode' do
      let(:params) { { q: package.conan_recipe } }

      subject { get api(url), params: params }
    end

    it_behaves_like 'conan search endpoint with access to package registry for everyone'
  end

  describe 'GET /api/v4/projects/:id/packages/conan/v2/conans/:package_name/:package_version/:package_username/' \
    ':package_channel/revisions/:recipe_revision/files/:file_name' do
    include_context 'for conan file download endpoints'

    let(:file_name) { recipe_file.file_name }
    let(:recipe_revision) { recipe_file_metadata.recipe_revision_value }
    let(:url_suffix) { "#{recipe_path}/revisions/#{recipe_revision}/files/#{file_name}" }

    subject(:request) { get api(url), headers: headers }

    it_behaves_like 'conan package revisions feature flag check'
    it_behaves_like 'packages feature check'
    it_behaves_like 'recipe file download endpoint'
    it_behaves_like 'accept get request on private project with access to package registry for everyone'
    it_behaves_like 'project not found by project id'

    it_behaves_like 'enforcing job token policies', :read_packages,
      allow_public_access_for_enabled_project_features: :package_registry do
      let(:headers) { job_basic_auth_header(target_job) }
    end

    it_behaves_like 'enforcing job token policies', :read_packages,
      allow_public_access_for_enabled_project_features: :package_registry do
      let(:headers) { job_basic_auth_header(target_job) }
    end

    describe 'parameter validation for recipe file endpoints' do
      using RSpec::Parameterized::TableSyntax

      let(:url_suffix) { "#{url_recipe_path}/revisions/#{url_recipe_revision}/files/#{url_file_name}" }

      # rubocop:disable Layout/LineLength -- Avoid formatting to keep one-line table syntax
      where(:error, :url_recipe_path, :url_recipe_revision, :url_file_name) do
        /package_name/     | 'pac$kage-1/1.0.0/namespace1+project-1/stable' | ref(:recipe_revision)                            | ref(:file_name)
        /package_version/  | 'package-1/1.0.$/namespace1+project-1/stable'  | ref(:recipe_revision)                            | ref(:file_name)
        /package_username/ | 'package-1/1.0.0/name$pace1+project-1/stable'  | ref(:recipe_revision)                            | ref(:file_name)
        /package_channel/  | 'package-1/1.0.0/namespace1+project-1/$table'  | ref(:recipe_revision)                            | ref(:file_name)
        /recipe_revision/  | ref(:recipe_path)                              | 'invalid_revi$ion'                               | ref(:file_name)
        /recipe_revision/  | ref(:recipe_path)                              | Packages::Conan::FileMetadatum::DEFAULT_REVISION | ref(:file_name)
        /file_name/        | ref(:recipe_path)                              | ref(:recipe_revision)                            | 'invalid_file.txt'
      end
      # rubocop:enable Layout/LineLength

      with_them do
        it_behaves_like 'returning response status with error', status: :bad_request, error: params[:error]
      end
    end
  end

  context 'with file upload endpoints' do
    include_context 'for conan file upload endpoints'
    let(:file_name) { 'conanfile.py' }
    let(:recipe_revision) { OpenSSL::Digest.hexdigest('MD5', 'valid_recipe_revision') }

    describe 'PUT /api/v4/projects/:id/packages/conan/v2/conans/:package_name/:package_version/:package_username/' \
      ':package_channel/revisions/:recipe_revision/files/:file_name' do
      let(:url_suffix) { "#{recipe_path}/revisions/#{recipe_revision}/files/#{file_name}" }

      subject(:request) { put api(url), headers: headers_with_token }

      it_behaves_like 'conan package revisions feature flag check'
      it_behaves_like 'packages feature check'
      it_behaves_like 'workhorse recipe file upload endpoint', recipe_revision: true
    end

    describe 'PUT /api/v4/projects/:id/packages/conan/v2/conans/:package_name/:package_version/:package_username/' \
      ':package_channel/revisions/:recipe_revision/files/:file_name/authorize' do
      let(:url_suffix) { "#{recipe_path}/revisions/#{recipe_revision}/files/#{file_name}/authorize" }

      subject(:request) do
        put api(url),
          headers: headers_with_token
      end

      it_behaves_like 'conan package revisions feature flag check'
      it_behaves_like 'packages feature check'
      it_behaves_like 'workhorse authorize endpoint'
    end
  end
end
