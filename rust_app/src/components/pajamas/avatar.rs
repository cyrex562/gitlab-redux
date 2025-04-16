use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// Avatar item types
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AvatarItem {
    /// User avatar
    User {
        id: i64,
        name: String,
        email: String,
        avatar_url: Option<String>,
        blocked: bool,
    },
    /// Project avatar
    Project {
        id: i64,
        name: String,
        avatar_url: Option<String>,
    },
    /// Group avatar
    Group {
        id: i64,
        name: String,
        avatar_url: Option<String>,
    },
    /// Email avatar
    Email(String),
    /// Direct URL avatar
    Url(String),
}

impl AvatarItem {
    fn name(&self) -> String {
        match self {
            AvatarItem::User { name, .. } => name.clone(),
            AvatarItem::Project { name, .. } => name.clone(),
            AvatarItem::Group { name, .. } => name.clone(),
            AvatarItem::Email(email) => email.clone(),
            AvatarItem::Url(url) => url.clone(),
        }
    }

    fn id(&self) -> Option<i64> {
        match self {
            AvatarItem::User { id, .. } => Some(*id),
            AvatarItem::Project { id, .. } => Some(*id),
            AvatarItem::Group { id, .. } => Some(*id),
            _ => None,
        }
    }

    fn is_user_or_email(&self) -> bool {
        matches!(self, AvatarItem::User { .. } | AvatarItem::Email(_))
    }

    fn avatar_url(&self, size: i32) -> Option<String> {
        match self {
            AvatarItem::User {
                avatar_url, email, ..
            } => avatar_url.clone().or_else(|| {
                Some(format!(
                    "https://secure.gravatar.com/avatar/{}?s={}&d=identicon",
                    md5::compute(email.as_bytes()),
                    size
                ))
            }),
            AvatarItem::Project { avatar_url, .. } | AvatarItem::Group { avatar_url, .. } => {
                avatar_url
                    .clone()
                    .map(|url| format!("{}?width={}", url, size))
            }
            AvatarItem::Email(email) => Some(format!(
                "https://secure.gravatar.com/avatar/{}?s={}&d=identicon",
                md5::compute(email.as_bytes()),
                size
            )),
            AvatarItem::Url(url) => Some(url.clone()),
        }
    }
}

/// Represents an email for avatar generation
#[derive(Debug, Clone)]
pub struct AvatarEmail {
    pub email: String,
}

impl AvatarEmail {
    pub fn new(email: String) -> Self {
        Self { email }
    }

    pub fn name(&self) -> String {
        self.email.clone()
    }
}

/// Size options for the avatar component
pub const SIZE_OPTIONS: [i32; 6] = [16, 24, 32, 48, 64, 96];

/// Renders a Pajamas compliant avatar element
#[derive(Debug, Clone)]
pub struct AvatarComponent {
    /// The item to display (User, Project, Group, AvatarEmail, or String URL)
    item: Box<dyn std::any::Any + Send + Sync>,
    /// Text for the alt attribute
    alt: Option<String>,
    /// Custom CSS class(es)
    class: String,
    /// Size in pixels
    size: i32,
    /// Additional avatar options
    avatar_options: HashMap<String, String>,
}

impl AvatarComponent {
    /// Create a new avatar component
    pub fn new<T: 'static + Send + Sync>(
        item: T,
        alt: Option<String>,
        class: String,
        size: i32,
        avatar_options: Option<HashMap<String, String>>,
    ) -> Self {
        Self {
            item: Box::new(item),
            alt,
            class,
            size: Self::filter_size(size),
            avatar_options: avatar_options.unwrap_or_default(),
        }
    }

    /// Filter size to ensure it's a valid option
    fn filter_size(size: i32) -> i32 {
        if SIZE_OPTIONS.contains(&size) {
            size
        } else {
            64 // default size
        }
    }

    /// Get the avatar classes
    fn avatar_classes(&self) -> String {
        let mut classes = vec![
            "gl-avatar".to_string(),
            format!("gl-avatar-s{}", self.size),
            self.class.clone(),
        ];

        // Add circle class for User or AvatarEmail
        if self.is_user_or_email() {
            classes.push("gl-avatar-circle".to_string());
        } else {
            classes.push("!gl-rounded-base".to_string());
        }

        // Add identicon classes if no src
        if self.src().is_none() {
            classes.push("gl-avatar-identicon".to_string());
            let id = self.get_id().unwrap_or(0);
            classes.push(format!("gl-avatar-identicon-bg{}", (id % 7) + 1));
        }

        classes.join(" ")
    }

    /// Check if item is User or AvatarEmail
    fn is_user_or_email(&self) -> bool {
        self.item.is::<User>() || self.item.is::<AvatarEmail>()
    }

    /// Get the item ID
    fn get_id(&self) -> Option<i64> {
        // This would need to be implemented based on your User/Project/Group types
        None
    }

    /// Get the source URL for the avatar
    fn src(&self) -> Option<String> {
        if let Some(url) = self.item.downcast_ref::<String>() {
            return Some(url.clone());
        }

        if let Some(user) = self.item.downcast_ref::<User>() {
            // This would need to be implemented based on your User type
            return None;
        }

        if let Some(email) = self.item.downcast_ref::<AvatarEmail>() {
            // This would need to be implemented based on your email avatar generation
            return None;
        }

        // Try to get avatar_url from the item
        if let Some(url) = self.get_avatar_url() {
            return Some(format!("{}?width={}", url, self.size));
        }

        None
    }

    /// Get the srcset for retina displays
    fn srcset(&self) -> Option<String> {
        self.src().map(|src| {
            let retina_src = src.replace(
                &format!("width={}", self.size),
                &format!("width={}", self.size * 2),
            );
            format!("{} 1x, {} 2x", src, retina_src)
        })
    }

    /// Get the alt text
    fn alt(&self) -> String {
        self.alt
            .clone()
            .unwrap_or_else(|| self.get_name().unwrap_or_default())
    }

    /// Get the initial for identicon
    fn initial(&self) -> String {
        self.get_name()
            .map(|name| {
                name.chars()
                    .next()
                    .unwrap_or(' ')
                    .to_uppercase()
                    .to_string()
            })
            .unwrap_or_default()
    }

    /// Get the name from the item
    fn get_name(&self) -> Option<String> {
        // This would need to be implemented based on your types
        None
    }

    /// Get the avatar URL from the item
    fn get_avatar_url(&self) -> Option<String> {
        // This would need to be implemented based on your types
        None
    }

    /// Render the avatar component
    pub fn render(&self) -> String {
        if let Some(src) = self.src() {
            let mut attrs = vec![
                format!("src=\"{}\"", src),
                format!("alt=\"{}\"", self.alt()),
                format!("class=\"{}\"", self.avatar_classes()),
                format!("height=\"{}\"", self.size),
                format!("width=\"{}\"", self.size),
                "loading=\"lazy\"".to_string(),
            ];

            if let Some(srcset) = self.srcset() {
                attrs.push(format!("srcset=\"{}\"", srcset));
            }

            // Add additional avatar options
            for (key, value) in &self.avatar_options {
                attrs.push(format!("{}=\"{}\"", key, value));
            }

            format!("<img {}>", attrs.join(" "))
        } else {
            let mut attrs = vec![
                format!("alt=\"{}\"", self.alt()),
                format!("class=\"{}\"", self.avatar_classes()),
            ];

            // Add additional avatar options
            for (key, value) in &self.avatar_options {
                attrs.push(format!("{}=\"{}\"", key, value));
            }

            format!("<div {}>{}</div>", attrs.join(" "), self.initial())
        }
    }
}

// This is a placeholder for the User type - you'll need to implement your own User type
#[derive(Debug, Clone)]
pub struct User {
    // Add necessary fields
}

// Helper trait for type checking
trait TypeCheck {
    fn is<T: 'static>(&self) -> bool;
}

impl TypeCheck for Box<dyn std::any::Any + Send + Sync> {
    fn is<T: 'static>(&self) -> bool {
        self.as_any().is::<T>()
    }
}

// Helper trait for downcasting
trait AsAny {
    fn as_any(&self) -> &dyn std::any::Any;
}

impl AsAny for Box<dyn std::any::Any + Send + Sync> {
    fn as_any(&self) -> &dyn std::any::Any {
        self.as_ref()
    }
}
