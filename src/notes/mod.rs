use actix_web::{web, HttpResponse, Responder};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::sync::Arc;

const MICROSECOND: i64 = 1_000_000;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Note {
    pub id: i64,
    pub note: String,
    pub noteable_type: String,
    pub noteable_id: i64,
    pub author_id: i64,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub position: Option<NotePosition>,
    pub confidential: bool,
    pub internal: bool,
    pub commit_id: Option<String>,
    pub discussion_id: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NotePosition {
    pub base_sha: String,
    pub start_sha: String,
    pub head_sha: String,
    pub position_type: String,
    pub old_path: Option<String>,
    pub new_path: Option<String>,
    pub old_line: Option<i32>,
    pub new_line: Option<i32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Discussion {
    pub id: String,
    pub notes: Vec<Note>,
    pub individual_note: bool,
    pub resolvable: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NoteCreate {
    pub note: String,
    pub position: Option<NotePosition>,
    pub confidential: Option<bool>,
    pub internal: Option<bool>,
    pub commit_id: Option<String>,
    pub merge_request_diff_head_sha: Option<String>,
    pub in_reply_to_discussion_id: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct NoteUpdate {
    pub note: String,
    pub position: Option<NotePosition>,
}

pub trait Noteable {
    fn id(&self) -> i64;
    fn type_name(&self) -> String;
}

pub struct NotesActionsHandler {
    db: Arc<dyn Database>,
    current_user: Option<User>,
    project: Option<Project>,
}

impl NotesActionsHandler {
    pub fn new(db: Arc<dyn Database>, current_user: Option<User>, project: Option<Project>) -> Self {
        Self { db, current_user, project }
    }

    pub async fn index(&self, noteable: &dyn Noteable, last_fetched_at: Option<i64>) -> impl Responder {
        let (notes, meta) = self.gather_all_notes(noteable, last_fetched_at).await;
        let notes = self.prepare_notes_for_rendering(notes);
        let notes = self.filter_readable_notes(notes);
        let notes = if self.use_note_serializer() {
            self.note_serializer().represent(notes)
        } else {
            notes.into_iter().map(|note| self.note_json(note)).collect()
        };

        HttpResponse::Ok().json(meta.merge(json!({ "notes": notes })))
    }

    pub async fn create(&self, noteable: &dyn Noteable, params: web::Json<NoteCreate>) -> impl Responder {
        let note = match self.create_note(noteable, params.into_inner()).await {
            Ok(note) => note,
            Err(e) => return HttpResponse::UnprocessableEntity().json(json!({
                "errors": e.to_string()
            }))
        };

        let mut json = json!({
            "commands_changes": note.commands_changes
        });

        if note.persisted() && self.return_discussion() {
            json["valid"] = json!(true);
            let discussion = note.discussion();
            let prepared_notes = self.prepare_notes_for_rendering(discussion.notes.clone());
            json["discussion"] = self.discussion_serializer().represent(discussion);
        } else {
            let prepared_notes = self.prepare_notes_for_rendering(vec![note.clone()]);
            json = json.merge(self.note_json(note));
        }

        if let Some(quick_actions) = note.quick_actions_status() {
            json["quick_actions_status"] = json!(quick_actions);
        }

        if note.has_errors() {
            HttpResponse::UnprocessableEntity().json(json!({
                "errors": note.errors().join(", ")
            }))
        } else if note.quick_actions_error() {
            HttpResponse::UnprocessableEntity().json(json!({
                "quick_actions_status": note.quick_actions_status()
            }))
        } else {
            HttpResponse::Ok().json(json)
        }
    }

    pub async fn update(&self, note: Note, params: web::Json<NoteUpdate>) -> impl Responder {
        let updated_note = match self.update_note(note, params.into_inner()).await {
            Ok(note) => note,
            Err(e) => return HttpResponse::UnprocessableEntity().json(json!({
                "errors": e.to_string()
            }))
        };

        if updated_note.is_destroyed() {
            return HttpResponse::Gone().finish();
        }

        let prepared_notes = self.prepare_notes_for_rendering(vec![updated_note.clone()]);
        HttpResponse::Ok().json(self.note_json(updated_note))
    }

    pub async fn destroy(&self, note: Note) -> impl Responder {
        if note.is_editable() {
            if let Err(e) = self.delete_note(note).await {
                return HttpResponse::InternalServerError().json(json!({
                    "message": e.to_string()
                }));
            }
        }
        HttpResponse::Ok().finish()
    }

    async fn gather_all_notes(&self, noteable: &dyn Noteable, last_fetched_at: Option<i64>) -> (Vec<Note>, serde_json::Value) {
        let now = Utc::now();
        let notes = self.merge_resource_events(
            self.notes_finder(noteable).execute().await,
            last_fetched_at
        ).await;

        let meta = json!({
            "last_fetched_at": (now.timestamp() * MICROSECOND) + (now.timestamp_subsec_micros() as i64)
        });

        (notes, meta)
    }

    async fn merge_resource_events(&self, notes: Vec<Note>, last_fetched_at: Option<i64>) -> Vec<Note> {
        if self.notes_filter() == NotesFilter::OnlyComments {
            return notes;
        }

        // TODO: Implement resource events merging
        notes
    }

    fn prepare_notes_for_rendering(&self, notes: Vec<Note>) -> Vec<Note> {
        // TODO: Implement note preparation
        notes
    }

    fn filter_readable_notes(&self, notes: Vec<Note>) -> Vec<Note> {
        notes.into_iter()
            .filter(|note| note.is_readable_by(self.current_user.as_ref()))
            .collect()
    }

    async fn create_note(&self, noteable: &dyn Noteable, params: NoteCreate) -> Result<Note, Box<dyn std::error::Error>> {
        // TODO: Implement note creation
        Err("Not implemented".into())
    }

    async fn update_note(&self, note: Note, params: NoteUpdate) -> Result<Note, Box<dyn std::error::Error>> {
        // TODO: Implement note update
        Err("Not implemented".into())
    }

    async fn delete_note(&self, note: Note) -> Result<(), Box<dyn std::error::Error>> {
        // TODO: Implement note deletion
        Ok(())
    }

    fn note_json(&self, note: Note) -> serde_json::Value {
        let mut attrs = json!({});

        if note.persisted() {
            attrs["valid"] = json!(true);

            if self.return_discussion() {
                let discussion = note.discussion();
                let prepared_notes = self.prepare_notes_for_rendering(discussion.notes.clone());
                attrs["discussion"] = self.discussion_serializer().represent(discussion);
            } else if self.use_note_serializer() {
                attrs = attrs.merge(self.note_serializer().represent(note));
            } else {
                attrs = attrs.merge(json!({
                    "id": note.id,
                    "discussion_id": note.discussion_id,
                    "html": self.note_html(note.clone()),
                    "note": note.note,
                    "on_image": note.is_on_image()
                }));

                let discussion = note.to_discussion();
                if !discussion.is_individual_note() {
                    attrs = attrs.merge(json!({
                        "discussion_resolvable": discussion.is_resolvable(),
                        "diff_discussion_html": self.diff_discussion_html(discussion.clone()),
                        "discussion_html": self.discussion_html(discussion.clone())
                    }));

                    if discussion.is_diff_discussion() {
                        attrs["discussion_line_code"] = json!(discussion.line_code());
                    }
                }
            }
        } else {
            attrs = attrs.merge(json!({
                "valid": false,
                "errors": note.errors()
            }));
        }

        attrs
    }

    fn note_html(&self, note: Note) -> String {
        // TODO: Implement note HTML rendering
        String::new()
    }

    fn diff_discussion_html(&self, discussion: Discussion) -> Option<String> {
        if !discussion.is_diff_discussion() {
            return None;
        }

        // TODO: Implement diff discussion HTML rendering
        None
    }

    fn discussion_html(&self, discussion: Discussion) -> Option<String> {
        if discussion.is_individual_note() {
            return None;
        }

        // TODO: Implement discussion HTML rendering
        None
    }

    fn notes_finder(&self, noteable: &dyn Noteable) -> NotesFinder {
        NotesFinder::new(self.current_user.clone(), noteable)
    }

    fn note_serializer(&self) -> NoteSerializer {
        NoteSerializer::new(self.project.clone(), noteable)
    }

    fn discussion_serializer(&self) -> DiscussionSerializer {
        DiscussionSerializer::new(self.current_user.clone())
    }

    fn use_note_serializer(&self) -> bool {
        !self.params.contains_key("html")
    }

    fn return_discussion(&self) -> bool {
        self.params.get("return_discussion")
            .and_then(|v| v.parse::<bool>().ok())
            .unwrap_or(false)
    }

    fn notes_filter(&self) -> NotesFilter {
        self.current_user
            .as_ref()
            .and_then(|user| user.notes_filter_for(self.params.get("target_type")))
            .unwrap_or(NotesFilter::All)
    }
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum NotesFilter {
    All,
    OnlyComments,
}

// These would be defined elsewhere in the codebase
pub trait Database {}
pub struct User {
    pub id: i64,
    // Add other user fields as needed
}
pub struct Project {
    pub id: i64,
    // Add other project fields as needed
}
pub struct NotesFinder {
    current_user: Option<User>,
    noteable: Box<dyn Noteable>,
}
pub struct NoteSerializer {
    project: Option<Project>,
    noteable: Box<dyn Noteable>,
}
pub struct DiscussionSerializer {
    current_user: Option<User>,
}

impl Note {
    fn persisted(&self) -> bool {
        self.id > 0
    }

    fn has_errors(&self) -> bool {
        // TODO: Implement error checking
        false
    }

    fn errors(&self) -> Vec<String> {
        // TODO: Implement error collection
        vec![]
    }

    fn is_editable(&self) -> bool {
        // TODO: Implement editability check
        true
    }

    fn is_readable_by(&self, user: Option<&User>) -> bool {
        // TODO: Implement readability check
        true
    }

    fn is_on_image(&self) -> bool {
        // TODO: Implement image check
        false
    }

    fn discussion(&self) -> Discussion {
        // TODO: Implement discussion retrieval
        Discussion {
            id: String::new(),
            notes: vec![],
            individual_note: true,
            resolvable: false,
        }
    }

    fn to_discussion(&self) -> Discussion {
        // TODO: Implement discussion conversion
        Discussion {
            id: String::new(),
            notes: vec![],
            individual_note: true,
            resolvable: false,
        }
    }

    fn commands_changes(&self) -> Option<serde_json::Value> {
        // TODO: Implement commands changes
        None
    }

    fn quick_actions_status(&self) -> Option<serde_json::Value> {
        // TODO: Implement quick actions status
        None
    }

    fn quick_actions_error(&self) -> bool {
        // TODO: Implement quick actions error check
        false
    }

    fn is_destroyed(&self) -> bool {
        // TODO: Implement destroyed check
        false
    }
}

impl Discussion {
    fn is_individual_note(&self) -> bool {
        self.individual_note
    }

    fn is_resolvable(&self) -> bool {
        self.resolvable
    }

    fn is_diff_discussion(&self) -> bool {
        // TODO: Implement diff discussion check
        false
    }

    fn line_code(&self) -> Option<String> {
        // TODO: Implement line code retrieval
        None
    }
}

impl NotesFinder {
    fn new(current_user: Option<User>, noteable: &dyn Noteable) -> Self {
        Self {
            current_user,
            noteable: Box::new(noteable.clone()),
        }
    }

    async fn execute(&self) -> Result<Vec<Note>, Box<dyn std::error::Error>> {
        // TODO: Implement notes finding
        Ok(vec![])
    }
}

impl NoteSerializer {
    fn new(project: Option<Project>, noteable: &dyn Noteable) -> Self {
        Self {
            project,
            noteable: Box::new(noteable.clone()),
        }
    }

    fn represent(&self, notes: Vec<Note>) -> serde_json::Value {
        // TODO: Implement note serialization
        json!({})
    }
}

impl DiscussionSerializer {
    fn new(current_user: Option<User>) -> Self {
        Self { current_user }
    }

    fn represent(&self, discussion: Discussion) -> serde_json::Value {
        // TODO: Implement discussion serialization
        json!({})
    }
} 