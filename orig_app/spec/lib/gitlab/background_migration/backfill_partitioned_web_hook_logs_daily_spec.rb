# frozen_string_literal: true

require 'spec_helper'

RSpec.describe Gitlab::BackgroundMigration::BackfillPartitionedWebHookLogsDaily,
  :migration_with_transaction,
  feature_category: :integrations do
  let(:connection) { ApplicationRecord.connection }
  let(:web_hook_logs) { table(:web_hook_logs, primary_key: :id) }
  let(:web_hook_logs_daily) { table(:web_hook_logs_daily, primary_key: :id) }
  let(:start_cursor) { [0, nil] }
  let(:end_cursor) { [web_hook_logs.last.id, Time.current.to_s] }
  let(:migration) do
    described_class.new(
      start_cursor: start_cursor,
      end_cursor: end_cursor,
      batch_table: :web_hook_logs,
      batch_column: :id,
      sub_batch_size: 1,
      pause_ms: 0,
      connection: connection
    )
  end

  before do
    connection.transaction do
      from = 1.month.ago.beginning_of_month
      to = 1.month.ago.end_of_month
      suffix = from.strftime('%Y%m')
      partition_name = "gitlab_partitions_dynamic.web_hook_logs_#{suffix}"

      current_month_start = Time.current.beginning_of_month
      current_month_end = current_month_start.end_of_month
      current_month_suffix = current_month_start.strftime('%Y%m')
      current_month_partition_name = "gitlab_partitions_dynamic.web_hook_logs_#{current_month_suffix}"

      connection.execute <<~SQL
        ALTER TABLE web_hook_logs DISABLE TRIGGER ALL; -- Don't sync records to partitioned table

        CREATE TABLE IF NOT EXISTS #{partition_name}
        PARTITION OF public.web_hook_logs
        FOR VALUES FROM (#{connection.quote(from)}) TO (#{connection.quote(to)});

        CREATE TABLE IF NOT EXISTS #{current_month_partition_name}
        PARTITION OF public.web_hook_logs
        FOR VALUES FROM (#{connection.quote(current_month_start)}) TO (#{connection.quote(current_month_end)});
      SQL

      create_web_hook_logs(created_at: from)
      create_web_hook_logs(created_at: 1.day.ago)

      connection.execute <<~SQL
        ALTER TABLE web_hook_logs ENABLE TRIGGER ALL;
      SQL
    end
  end

  describe '#perform' do
    it 'backfills web_hook_logs_daily from web_hook_logs only for existing partition' do
      migration.perform

      expect(web_hook_logs_daily.count).to eq(1)
    end
  end

  private

  def create_web_hook_logs(**params)
    web_hook_logs_params = {
      web_hook_id: 1,
      trigger: 'push',
      url: 'https://example.com/webhook',
      request_headers: { "Content-Type": "application/json" },
      request_data: { key: "value" },
      response_headers: { Server: "nginx" },
      response_body: { status: "success" },
      response_status: '200',
      execution_duration: 0.5,
      url_hash: 'abc123',
      updated_at: params[:created_at]
    }

    web_hook_logs_params.merge!(params)

    web_hook_logs.create!(web_hook_logs_params)
  end
end
