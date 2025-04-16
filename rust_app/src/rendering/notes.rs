use super::{RenderOptions, Renderable};
use actix_web::{web, HttpResponse};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Note {
    pub id: i32,
    pub body: String,
    pub author_id: i32,
    pub author_name: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub position: Option<NotePosition>,
    pub system: bool,
    pub noteable_id: i32,
    pub noteable_type: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct NotePosition {
    pub base_sha: String,
    pub start_sha: String,
    pub head_sha: String,
    pub old_path: Option<String>,
    pub new_path: Option<String>,
    pub position_type: String,
    pub old_line: Option<i32>,
    pub new_line: Option<i32>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct NoteCollection {
    pub notes: Vec<Note>,
    pub total_count: usize,
}

impl Note {
    pub fn new(
        id: i32,
        body: String,
        author_id: i32,
        author_name: String,
        created_at: DateTime<Utc>,
        updated_at: DateTime<Utc>,
        noteable_id: i32,
        noteable_type: String,
    ) -> Self {
        Self {
            id,
            body,
            author_id,
            author_name,
            created_at,
            updated_at,
            position: None,
            system: false,
            noteable_id,
            noteable_type,
        }
    }

    pub fn with_position(mut self, position: NotePosition) -> Self {
        self.position = Some(position);
        self
    }

    pub fn with_system(mut self, system: bool) -> Self {
        self.system = system;
        self
    }
}

impl Renderable for Note {
    fn render(&self, options: &RenderOptions) -> HttpResponse {
        // TODO: Implement note rendering
        HttpResponse::Ok().json(self)
    }

    fn render_json(&self, options: &RenderOptions) -> HttpResponse {
        self.render(options)
    }
}

impl Renderable for NoteCollection {
    fn render(&self, options: &RenderOptions) -> HttpResponse {
        // TODO: Implement note collection rendering
        HttpResponse::Ok().json(self)
    }

    fn render_json(&self, options: &RenderOptions) -> HttpResponse {
        self.render(options)
    }
}

pub struct NoteRenderer;

impl NoteRenderer {
    pub fn new() -> Self {
        Self
    }

    pub fn render_note(&self, note: &Note, options: &RenderOptions) -> HttpResponse {
        note.render(options)
    }

    pub fn render_notes(&self, notes: &[Note], options: &RenderOptions) -> HttpResponse {
        let collection = NoteCollection {
            notes: notes.to_vec(),
            total_count: notes.len(),
        };
        collection.render(options)
    }

    pub fn prepare_notes_for_rendering(&self, notes: &[Note]) -> Vec<Note> {
        // TODO: Implement note preparation for rendering
        notes.to_vec()
    }
}
