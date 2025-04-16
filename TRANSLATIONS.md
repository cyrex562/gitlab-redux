## Form Components

### Base Components

- `orig_app/app/components/pajamas/component.rb` -> `rust_app/src/components/pajamas/component.rs`
- `orig_app/app/components/pajamas/concerns/checkbox_radio_label_with_help_text.rb` -> `rust_app/src/components/pajamas/concerns/checkbox_radio_label_with_help_text.rs`
- `orig_app/app/components/pajamas/concerns/checkbox_radio_options.rb` -> `rust_app/src/components/pajamas/concerns/checkbox_radio_options.rs`

### Checkbox Components

- `orig_app/app/components/pajamas/checkbox_component.rb` -> `rust_app/src/components/pajamas/checkbox.rs`
- `orig_app/app/components/pajamas/checkbox_component.html.haml` -> Integrated into `checkbox.rs`
- `orig_app/app/components/pajamas/checkbox_tag_component.rb` -> `rust_app/src/components/pajamas/checkbox_tag.rs`
- `orig_app/app/components/pajamas/checkbox_tag_component.html.haml` -> Integrated into `checkbox_tag.rs`

### Radio Components

- `orig_app/app/components/pajamas/radio_component.rb` -> `rust_app/src/components/pajamas/radio.rs`
- `orig_app/app/components/pajamas/radio_component.html.haml` -> Integrated into `radio.rs`
- `orig_app/app/components/pajamas/radio_tag_component.rb` -> `rust_app/src/components/pajamas/radio_tag.rs`
- `orig_app/app/components/pajamas/radio_tag_component.html.haml` -> Integrated into `radio_tag.rs`

### Toggle Component

- `orig_app/app/components/pajamas/toggle_component.rb` -> `rust_app/src/components/pajamas/toggle.rs`
- `orig_app/app/components/pajamas/toggle_component.html.haml` -> Integrated into `toggle.rs`

## Controllers

- `orig_app/app/controllers/application_controller.rb` -> `rust_app/src/controllers/mod.rs`
- `orig_app/app/controllers/users_controller.rb` -> `rust_app/src/controllers/users.rs`
- `orig_app/app/controllers/sessions_controller.rb` -> `rust_app/src/controllers/sessions.rs`
