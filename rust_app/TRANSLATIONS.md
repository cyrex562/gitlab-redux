# Ruby to Rust Translations

This file documents the translations from Ruby to Rust in this project.

## Application Cable

- `orig_app/app/channels/application_cable/channel.rb` -> `rust_app/src/websocket/channel.rs`
- `orig_app/app/channels/application_cable/connection.rb` -> `rust_app/src/websocket/connection.rs`
- `orig_app/app/channels/application_cable/logging.rb` -> `rust_app/src/websocket/logging.rs`

## Noteable

- `orig_app/app/channels/noteable/notes_channel.rb` -> `rust_app/src/websocket/noteable/notes_channel.rs`

## Diffs Components

- `orig_app/app/components/diffs/base_component.rb` -> `rust_app/src/components/diffs/base.rs`
- `orig_app/app/components/diffs/overflow_warning_component.rb` -> `rust_app/src/components/diffs/overflow_warning.rs`
- `orig_app/app/components/diffs/overflow_warning_component.html.haml` -> `rust_app/src/components/diffs/overflow_warning.rs` (HTML template)
- `orig_app/app/components/diffs/stats_component.rb` -> `rust_app/src/components/diffs/stats.rs`
- `orig_app/app/components/diffs/stats_component.html.haml` -> `rust_app/src/components/diffs/stats.rs` (HTML template)
