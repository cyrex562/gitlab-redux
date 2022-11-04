# frozen_string_literal: true

module Gitlab
  module Database
    DATABASE_NAMES = %w[main ci].freeze

    MAIN_DATABASE_NAME = 'main'
    CI_DATABASE_NAME = 'ci'
    DEFAULT_POOL_HEADROOM = 10

    # This constant is used when renaming tables concurrently.
    # If you plan to rename a table using the `rename_table_safely` method, add your table here one milestone before the rename.
    # Example:
    # TABLES_TO_BE_RENAMED = {
    #   'old_name' => 'new_name'
    # }.freeze
    TABLES_TO_BE_RENAMED = {}.freeze

    # Minimum PostgreSQL version requirement per documentation:
    # https://docs.gitlab.com/ee/install/requirements.html#postgresql-requirements
    MINIMUM_POSTGRES_VERSION = 12

    # https://www.postgresql.org/docs/9.2/static/datatype-numeric.html
    MAX_INT_VALUE = 2147483647
    MIN_INT_VALUE = -2147483648

    # The max value between MySQL's TIMESTAMP and PostgreSQL's timestampz:
    # https://www.postgresql.org/docs/9.1/static/datatype-datetime.html
    # https://dev.mysql.com/doc/refman/5.7/en/datetime.html
    # FIXME: this should just be the max value of timestampz
    MAX_TIMESTAMP_VALUE = Time.at((1 << 31) - 1).freeze

    # The maximum number of characters for text fields, to avoid DoS attacks via parsing huge text fields
    # https://gitlab.com/gitlab-org/gitlab-foss/issues/61974
    MAX_TEXT_SIZE_LIMIT = 1_000_000

    # Minimum schema version from which migrations are supported
    # Migrations before this version may have been removed
    MIN_SCHEMA_VERSION = 20190506135400
    MIN_SCHEMA_GITLAB_VERSION = '11.11.0'

    # Schema we store dynamically managed partitions in (e.g. for time partitioning)
    DYNAMIC_PARTITIONS_SCHEMA = :gitlab_partitions_dynamic

    # Schema we store static partitions in (e.g. for hash partitioning)
    STATIC_PARTITIONS_SCHEMA = :gitlab_partitions_static

    # This is an extensive list of postgres schemas owned by GitLab
    # It does not include the default public schema
    EXTRA_SCHEMAS = [DYNAMIC_PARTITIONS_SCHEMA, STATIC_PARTITIONS_SCHEMA].freeze

    PRIMARY_DATABASE_NAME = ActiveRecord::Base.connection_db_config.name.to_sym # rubocop:disable Database/MultipleDatabases

    def self.database_base_models
      @database_base_models ||= {
        # Note that we use ActiveRecord::Base here and not ApplicationRecord.
        # This is deliberate, as we also use these classes to apply load
        # balancing to, and the load balancer must be enabled for _all_ models
        # that inherit from ActiveRecord::Base; not just our own models that
        # inherit from ApplicationRecord.
        main: ::ActiveRecord::Base,
        ci: ::Ci::ApplicationRecord.connection_class? ? ::Ci::ApplicationRecord : nil
      }.compact.with_indifferent_access.freeze
    end

    # This returns a list of databases that contains all the gitlab_shared schema
    # tables. We can't reuse database_base_models because Geo does not support
    # the gitlab_shared tables yet.
    def self.database_base_models_with_gitlab_shared
      @database_base_models_with_gitlab_shared ||= {
        # Note that we use ActiveRecord::Base here and not ApplicationRecord.
        # This is deliberate, as we also use these classes to apply load
        # balancing to, and the load balancer must be enabled for _all_ models
        # that inher from ActiveRecord::Base; not just our own models that
        # inherit from ApplicationRecord.
        main: ::ActiveRecord::Base,
        ci: ::Ci::ApplicationRecord.connection_class? ? ::Ci::ApplicationRecord : nil
      }.compact.with_indifferent_access.freeze
    end

    # This returns a list of databases whose connection supports database load
    # balancing. We can't reuse the database_base_models method because the Geo
    # database does not support load balancing yet.
    #
    # TODO: https://gitlab.com/gitlab-org/geo-team/discussions/-/issues/5032
    def self.database_base_models_using_load_balancing
      @database_base_models_with_gitlab_shared ||= {
        # Note that we use ActiveRecord::Base here and not ApplicationRecord.
        # This is deliberate, as we also use these classes to apply load
        # balancing to, and the load balancer must be enabled for _all_ models
        # that inher from ActiveRecord::Base; not just our own models that
        # inherit from ApplicationRecord.
        main: ::ActiveRecord::Base,
        ci: ::Ci::ApplicationRecord.connection_class? ? ::Ci::ApplicationRecord : nil
      }.compact.with_indifferent_access.freeze
    end

    # This returns a list of base models with connection associated for a given gitlab_schema
    def self.schemas_to_base_models
      @schemas_to_base_models ||= {
        gitlab_main: [self.database_base_models.fetch(:main)],
        gitlab_ci: [self.database_base_models[:ci] || self.database_base_models.fetch(:main)], # use CI or fallback to main
        gitlab_shared: database_base_models_with_gitlab_shared.values, # all models
        gitlab_internal: database_base_models.values # all models
      }.with_indifferent_access.freeze
    end

    def self.all_database_names
      DATABASE_NAMES
    end

    # We configure the database connection pool size automatically based on the
    # configured concurrency. We also add some headroom, to make sure we don't
    # run out of connections when more threads besides the 'user-facing' ones
    # are running.
    #
    # Read more about this in
    # doc/development/database/client_side_connection_pool.md
    def self.default_pool_size
      headroom =
        (ENV["DB_POOL_HEADROOM"].presence || DEFAULT_POOL_HEADROOM).to_i

      Gitlab::Runtime.max_threads + headroom
    end

    def self.has_config?(database_name)
      Gitlab::Application.config.database_configuration[Rails.env].include?(database_name.to_s)
    end

    class PgUser < ApplicationRecord
      self.table_name = 'pg_user'
      self.primary_key = :usename
    end

    # rubocop: disable CodeReuse/ActiveRecord
    def self.check_for_non_superuser
      user = PgUser.find_by('usename = CURRENT_USER')
      am_i_superuser = user.usesuper

      Gitlab::AppLogger.info(
        "Account details: User: \"#{user.usename}\", UseSuper: (#{am_i_superuser})"
      )

      raise 'Error: detected superuser' if am_i_superuser
    rescue ActiveRecord::StatementInvalid
      raise 'User CURRENT_USER not found'
    end
    # rubocop: enable CodeReuse/ActiveRecord

    def self.check_postgres_version_and_print_warning
      return if Gitlab::Runtime.rails_runner?

      database_base_models.each do |name, model|
        database = Gitlab::Database::Reflection.new(model)

        next if database.postgresql_minimum_supported_version?

        Kernel.warn ERB.new(Rainbow.new.wrap(<<~EOS).red).result

                    ██     ██  █████  ██████  ███    ██ ██ ███    ██  ██████ 
                    ██     ██ ██   ██ ██   ██ ████   ██ ██ ████   ██ ██      
                    ██  █  ██ ███████ ██████  ██ ██  ██ ██ ██ ██  ██ ██   ███ 
                    ██ ███ ██ ██   ██ ██   ██ ██  ██ ██ ██ ██  ██ ██ ██    ██ 
                     ███ ███  ██   ██ ██   ██ ██   ████ ██ ██   ████  ██████  

          ******************************************************************************
            You are using PostgreSQL #{database.version} for the #{name} database, but PostgreSQL >= <%= Gitlab::Database::MINIMUM_POSTGRES_VERSION %>
            is required for this version of GitLab.
            <% if Rails.env.development? || Rails.env.test? %>
            If using gitlab-development-kit, please find the relevant steps here:
              https://gitlab.com/gitlab-org/gitlab-development-kit/-/blob/main/doc/howto/postgresql.md#upgrade-postgresql
            <% end %>
            Please upgrade your environment to a supported PostgreSQL version, see
            https://docs.gitlab.com/ee/install/requirements.html#database for details.
          ******************************************************************************
        EOS
      rescue ActiveRecord::ActiveRecordError, PG::Error
        # ignore - happens when Rake tasks yet have to create a database, e.g. for testing
      end
    end

    def self.random
      "RANDOM()"
    end

    def self.true_value
      "'t'"
    end

    def self.false_value
      "'f'"
    end

    def self.sanitize_timestamp(timestamp)
      MAX_TIMESTAMP_VALUE > timestamp ? timestamp : MAX_TIMESTAMP_VALUE.dup
    end

    def self.all_uncached(&block)
      # Calls to #uncached only disable caching for the current connection. Since the load balancer
      # can potentially upgrade from read to read-write mode (using a different connection), we specify
      # up-front that we'll explicitly use the primary for the duration of the operation.
      Gitlab::Database::LoadBalancing::Session.current.use_primary do
        base_models = database_base_models_using_load_balancing.values
        base_models.reduce(block) { |blk, model| -> { model.uncached(&blk) } }.call
      end
    end

    def self.allow_cross_joins_across_databases(url:)
      # this method is implemented in:
      # spec/support/database/prevent_cross_joins.rb
      yield
    end

    def self.add_post_migrate_path_to_rails(force: false)
      return if ENV['SKIP_POST_DEPLOYMENT_MIGRATIONS'] && !force

      Rails.application.config.paths['db'].each do |db_path|
        path = Rails.root.join(db_path, 'post_migrate').to_s

        unless Rails.application.config.paths['db/migrate'].include? path
          Rails.application.config.paths['db/migrate'] << path

          # Rails memoizes migrations at certain points where it won't read the above
          # path just yet. As such we must also update the following list of paths.
          ActiveRecord::Migrator.migrations_paths << path
        end
      end
    end

    def self.db_config_names
      ::ActiveRecord::Base.configurations.configs_for(env_name: Rails.env).map(&:name) - ['geo']
    end

    # This returns all matching schemas that a given connection can use
    # Since the `ActiveRecord::Base` might change the connection (from main to ci)
    # This does not look at literal connection names, but rather compares
    # models that are holders for a given db_config_name
    def self.gitlab_schemas_for_connection(connection)
      db_config = self.db_config_for_connection(connection)

      # connection might not be yet adopted (returning NullPool, and no connection_klass)
      # in such cases it is fine to ignore such connections
      return unless db_config

      db_config_name = db_config.name.delete_suffix(LoadBalancing::LoadBalancer::REPLICA_SUFFIX)
      primary_model = self.database_base_models.fetch(db_config_name.to_sym)

      self.schemas_to_base_models.select do |_, child_models|
        child_models.any? do |child_model|
          child_model == primary_model || \
            # The model might indicate a child connection, ensure that this is enclosed in a `db_config`
            self.database_base_models[self.db_config_share_with(child_model.connection_db_config)] == primary_model
        end
      end.keys.map!(&:to_sym)
    end

    def self.db_config_for_connection(connection)
      return unless connection

      # For a ConnectionProxy we want to avoid ambiguous db_config as it may
      # sometimes default to replica so we always return the primary config
      # instead.
      if connection.is_a?(::Gitlab::Database::LoadBalancing::ConnectionProxy)
        return connection.load_balancer.configuration.db_config
      end

      # During application init we might receive `NullPool`
      return unless connection.respond_to?(:pool) &&
        connection.pool.respond_to?(:db_config)

      connection.pool.db_config
    end

    # At the moment, the connection can only be retrieved by
    # Gitlab::Database::LoadBalancer#read or #read_write or from the
    # ActiveRecord directly. Therefore, if the load balancer doesn't
    # recognize the connection, this method returns the primary role
    # directly. In future, we may need to check for other sources.
    # Expected returned names:
    # main, main_replica, ci, ci_replica, unknown
    def self.db_config_name(connection)
      db_config = db_config_for_connection(connection)
      db_config&.name || 'unknown'
    end

    # Currently the database configuration can only be shared with `main:`
    # If the `database_tasks: false` is being used
    # This is to be refined: https://gitlab.com/gitlab-org/gitlab/-/issues/356580
    def self.db_config_share_with(db_config)
      if db_config.database_tasks?
        nil # no sharing
      else
        'main' # share with `main:`
      end
    end

    def self.read_only?
      false
    end

    def self.read_write?
      !read_only?
    end

    # Monkeypatch rails with upgraded database observability
    def self.install_transaction_metrics_patches!
      ActiveRecord::Base.prepend(ActiveRecordBaseTransactionMetrics)
    end

    def self.install_transaction_context_patches!
      ActiveRecord::ConnectionAdapters::TransactionManager
        .prepend(TransactionManagerContext)
      ActiveRecord::ConnectionAdapters::RealTransaction
        .prepend(RealTransactionContext)
    end

    # MonkeyPatch for ActiveRecord::Base for adding observability
    module ActiveRecordBaseTransactionMetrics
      extend ActiveSupport::Concern

      class_methods do
        # A patch over ApplicationRecord.transaction that provides
        # observability into transactional methods.
        def transaction(**options, &block)
          transaction_type = get_transaction_type(connection.transaction_open?, options[:requires_new])

          ::Gitlab::Database::Metrics.subtransactions_increment(self.name) if transaction_type == :sub_transaction

          payload = { connection: connection, transaction_type: transaction_type }

          ActiveSupport::Notifications.instrument('transaction.active_record', payload) do
            super(**options, &block)
          end
        end

        private

        def get_transaction_type(transaction_open, requires_new_flag)
          if transaction_open
            return :sub_transaction if requires_new_flag

            return :fake_transaction
          end

          :real_transaction
        end
      end
    end

    # rubocop:disable Gitlab/ModuleWithInstanceVariables
    module TransactionManagerContext
      def transaction_context
        @stack.first.try(:gitlab_transaction_context)
      end
    end

    module RealTransactionContext
      def gitlab_transaction_context
        @gitlab_transaction_context ||= ::Gitlab::Database::Transaction::Context.new
      end

      def commit
        gitlab_transaction_context.commit

        super
      end

      def rollback
        gitlab_transaction_context.rollback

        super
      end
    end
    # rubocop:enable Gitlab/ModuleWithInstanceVariables
  end
end

Gitlab::Database.prepend_mod_with('Gitlab::Database')
