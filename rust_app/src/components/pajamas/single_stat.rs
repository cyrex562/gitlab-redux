use std::collections::HashMap;

use super::badge::{BadgeComponent, BadgeVariant};

/// Renders a Pajamas compliant single stat component
#[derive(Debug, Clone)]
pub struct SingleStatComponent {
    /// The title of the stat
    title: Option<String>,
    /// The title tag (e.g., 'span', 'h1', etc.)
    title_tag: String,
    /// The stat value
    stat_value: Option<String>,
    /// The test ID for the stat value
    stat_value_testid: String,
    /// The unit for the stat value
    unit: Option<String>,
    /// The title icon
    title_icon: Option<String>,
    /// The meta text
    meta_text: Option<String>,
    /// The meta icon
    meta_icon: Option<String>,
    /// The text color
    text_color: Option<String>,
    /// The variant of the badge
    variant: BadgeVariant,
    /// The title content slot
    title_slot: Option<String>,
    /// The stat value content slot
    stat_value_slot: Option<String>,
}

impl SingleStatComponent {
    /// Create a new single stat component
    pub fn new(
        title: Option<String>,
        title_tag: Option<String>,
        stat_value: Option<String>,
        stat_value_testid: Option<String>,
        unit: Option<String>,
        title_icon: Option<String>,
        meta_text: Option<String>,
        meta_icon: Option<String>,
        text_color: Option<String>,
        variant: Option<BadgeVariant>,
    ) -> Self {
        Self {
            title,
            title_tag: title_tag.unwrap_or_else(|| "span".to_string()),
            stat_value,
            stat_value_testid: stat_value_testid
                .unwrap_or_else(|| "non-animated-value".to_string()),
            unit,
            title_icon: title_icon.filter(|s| !s.is_empty()),
            meta_text,
            meta_icon,
            text_color,
            variant: variant.unwrap_or(BadgeVariant::Muted),
            title_slot: None,
            stat_value_slot: None,
        }
    }

    /// Set the title content slot
    pub fn with_title_slot(mut self, title: String) -> Self {
        self.title_slot = Some(title);
        self
    }

    /// Set the stat value content slot
    pub fn with_stat_value_slot(mut self, stat_value: String) -> Self {
        self.stat_value_slot = Some(stat_value);
        self
    }

    /// Get the title content
    fn get_title(&self) -> Option<String> {
        self.title_slot.clone().or_else(|| self.title.clone())
    }

    /// Get the stat value content
    fn get_stat_value(&self) -> Option<String> {
        self.stat_value_slot
            .clone()
            .or_else(|| self.stat_value.clone())
    }

    /// Check if the unit is present
    fn has_unit(&self) -> bool {
        self.unit.is_some()
    }

    /// Check if the title icon is present
    fn has_title_icon(&self) -> bool {
        self.title_icon.is_some()
    }

    /// Check if the meta icon is present
    fn has_meta_icon(&self) -> bool {
        self.meta_icon.is_some()
    }

    /// Check if the meta text is present
    fn has_meta_text(&self) -> bool {
        self.meta_text.is_some()
    }

    /// Get the unit class
    fn unit_class(&self) -> String {
        if !self.has_unit() {
            "gl-mr-2".to_string()
        } else {
            String::new()
        }
    }

    /// Render the title icon
    fn render_title_icon(&self) -> Option<String> {
        self.title_icon.as_ref().map(|icon| {
            format!(
                r#"<svg class="s16 gl-mr-2 gl-fill-icon-subtle" data-testid="{}"><use xlink:href="/assets/icons.svg#{}"></use></svg>"#,
                icon, icon
            )
        })
    }

    /// Render the meta icon
    fn render_meta_icon(&self) -> Option<String> {
        if let Some(icon) = &self.meta_icon {
            let text_color = self
                .text_color
                .as_ref()
                .map(|c| format!(" class=\"{}\"", c))
                .unwrap_or_default();

            Some(format!(
                r#"<svg class="s16{}" data-testid="{}"><use xlink:href="/assets/icons.svg#{}"></use></svg>"#,
                text_color, icon, icon
            ))
        } else {
            None
        }
    }

    /// Render the meta badge
    fn render_meta_badge(&self) -> Option<String> {
        if let Some(text) = &self.meta_text {
            let mut badge = BadgeComponent::new(text.clone(), self.variant);

            if let Some(icon) = &self.meta_icon {
                badge = badge.with_icon(icon.clone());
            }

            Some(badge.render())
        } else {
            None
        }
    }

    /// Render the single stat component
    pub fn render(&self) -> String {
        let mut html = String::new();

        // Single stat container
        html.push_str(r#"<div class="gl-single-stat gl-flex gl-flex-col gl-py-2">"#);

        // Title section
        html.push_str(r#"<div class="gl-flex gl-items-center gl-text-subtle gl-mb-4">"#);

        // Title icon
        if let Some(title_icon) = self.render_title_icon() {
            html.push_str(&title_icon);
        }

        // Title
        if let Some(title) = self.get_title() {
            html.push_str(&format!(
                r#"<{} class="gl-text-base gl-font-normal gl-text-subtle gl-leading-reset gl-m-0" data-testid="title-text">{}</{}>"#,
                self.title_tag, title, self.title_tag
            ));
        }

        html.push_str("</div>");

        // Stat value section
        html.push_str(r#"<div class="gl-single-stat-content gl-flex gl-items-baseline gl-font-bold gl-text-default gl-mb-4">"#);

        // Stat value
        if let Some(stat_value) = self.get_stat_value() {
            let unit_class = self.unit_class();
            html.push_str(&format!(
                r#"<span class="gl-single-stat-number gl-leading-1{}" data-testid="displayValue"><span data-testid="{}">{}</span></span>"#,
                unit_class, self.stat_value_testid, stat_value
            ));
        }

        // Unit
        if let Some(unit) = &self.unit {
            html.push_str(&format!(
                r#"<span class="gl-text-sm gl-mx-2 gl-transition-all gl-opacity-10" data-testid="unit">{}</span>"#,
                unit
            ));
        }

        // Meta icon or badge
        if self.has_meta_icon() && !self.has_meta_text() {
            if let Some(meta_icon) = self.render_meta_icon() {
                html.push_str(&meta_icon);
            }
        } else if self.has_meta_text() {
            if let Some(meta_badge) = self.render_meta_badge() {
                html.push_str(&meta_badge);
            }
        }

        html.push_str("</div></div>");
        html
    }
}
