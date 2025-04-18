use crate::models::note::Note;
use crate::models::project::Project;
use crate::models::user::User;
use crate::services::notes::{CreateService, DestroyService, UpdateService};
use crate::services::notes_finder::NotesFinder;
use crate::services::resource_events::MergeIntoNotesService;
use crate::utils::boolean::to_boolean;
use crate::utils::polling_interval::PollingInterval;
use crate::utils::strong_memoize::StrongMemoize;
use actix_web::{web, HttpRequest, HttpResponse, Responder};
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

const MICROSECOND: i64 = 1_000_000;

#[derive(Debug, Serialize, Deserialize)]
pub struct CreateNoteParams {
    pub note: String,
    pub position: Option<Position>,
    pub confidential: Option<bool>,
    pub internal: Option<bool>,
    pub merge_request_diff_head_sha: Option<String>,
    pub in_reply_to_discussion_id: Option<String>,
    pub commit_id: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Position {
    pub base_sha: String,
    pub start_sha: String,
    pub head_sha: String,
    pub position_type: String,
    pub old_path: Option<String>,
    pub new_path: Option<String>,
    pub old_line: Option<i32>,
    pub new_line: Option<i32>,
}

pub trait NotesActions {
    fn index(&self, req: HttpRequest) -> impl Responder;
    fn create(&self, req: HttpRequest) -> impl Responder;
    fn update(&self, req: HttpRequest) -> impl Responder;
    fn destroy(&self, req: HttpRequest) -> impl Responder;
}

pub trait Noteable {
    fn id(&self) -> i32;
    fn class_name(&self) -> String;
}

pub struct NotesActionsImpl {
    current_user: User,
    project: Project,
    noteable: Option<Box<dyn Noteable>>,
    note: Option<Note>,
}

impl NotesActionsImpl {
    pub fn new(
        current_user: User,
        project: Project,
        noteable: Option<Box<dyn Noteable>>,
        note: Option<Note>,
    ) -> Self {
        Self {
            current_user,
            project,
            noteable,
            note,
        }
    }

    fn gather_all_notes(&self) -> (Vec<Note>, serde_json::Value) {
        let now = Utc::now();
        let notes = self.merge_resource_events(
            self.notes_finder()
                .execute()
                .inc_relations_for_view(self.noteable()),
        );

        let meta = serde_json::json!({
            "last_fetched_at": (now.timestamp() * MICROSECOND) + (now.timestamp_subsec_micros() as i64)
        });

        (notes, meta)
    }

    fn merge_resource_events(&self, notes: Vec<Note>) -> Vec<Note> {
        if self.notes_filter() == "only_comments" {
            return notes;
        }

        MergeIntoNotesService::new(self.noteable(), &self.current_user, self.last_fetched_at())
            .execute(notes)
    }

    fn note_html(&self, note: &Note) -> String {
        // Implementation depends on your template engine
        // This is a placeholder that should be replaced with actual template rendering
        format!("Rendered note: {:?}", note)
    }

    fn note_json(&self, note: &Note) -> serde_json::Value {
        let mut attrs = serde_json::json!({});

        if note.is_persisted() {
            attrs["valid"] = serde_json::json!(true);

            if self.return_discussion() {
                let discussion = note.discussion();
                self.prepare_notes_for_rendering(discussion.notes());

                attrs["discussion"] = self.discussion_serializer().represent(discussion, self);
            } else if self.use_note_serializer() {
                attrs = self.note_serializer().represent(note);
            } else {
                attrs = serde_json::json!({
                    "id": note.id,
                    "discussion_id": note.discussion_id(self.noteable()),
                    "html": self.note_html(note),
                    "note": note.note,
                    "on_image": note.on_image()
                });

                let discussion = note.to_discussion(self.noteable());
                if !discussion.individual_note() {
                    attrs["discussion_resolvable"] = serde_json::json!(discussion.resolvable());
                    attrs["diff_discussion_html"] =
                        serde_json::json!(self.diff_discussion_html(&discussion));
                    attrs["discussion_html"] = serde_json::json!(self.discussion_html(&discussion));

                    if discussion.diff_discussion() {
                        attrs["discussion_line_code"] = serde_json::json!(discussion.line_code());
                    }
                }
            }
        } else {
            attrs = serde_json::json!({
                "valid": false,
                "errors": note.errors()
            });
        }

        attrs
    }

    fn diff_discussion_html(&self, discussion: &Discussion) -> Option<String> {
        if !discussion.diff_discussion() {
            return None;
        }

        let on_image = discussion.on_image();
        let template = if self.params().get("view") == Some("parallel") && !on_image {
            "discussions/_parallel_diff_discussion"
        } else {
            "discussions/_diff_discussion"
        };

        let locals = if self.params().get("line_type") == Some("old") {
            serde_json::json!({
                "discussions_left": vec![discussion],
                "discussions_right": null
            })
        } else {
            serde_json::json!({
                "discussions_left": null,
                "discussions_right": vec![discussion]
            })
        };

        Some(self.view_to_html_string(template, locals))
    }

    fn discussion_html(&self, discussion: &Discussion) -> Option<String> {
        if discussion.individual_note() {
            return None;
        }

        Some(self.view_to_html_string(
            "discussions/_discussion",
            serde_json::json!({ "discussion": discussion }),
        ))
    }

    fn prepare_notes_for_rendering(&self, notes: Vec<Note>) {
        // Implementation depends on your note preparation logic
    }

    fn notes_filter(&self) -> String {
        self.current_user
            .notes_filter_for(self.params().get("target_type"))
    }

    fn notes_finder(&self) -> NotesFinder {
        NotesFinder::new(&this.current_user, self.finder_params())
    }

    fn finder_params(&self) -> serde_json::Value {
        serde_json::json!({
            "target_type": self.noteable().map(|n| n.class_name()),
            "target_id": self.noteable().map(|n| n.id()),
            "last_fetched_at": self.last_fetched_at()
        })
    }

    fn last_fetched_at(&self) -> Option<i64> {
        self.req
            .headers()
            .get("X-Last-Fetched-At")
            .and_then(|h| h.to_str().ok())
            .and_then(|s| s.parse::<i64>().ok())
    }

    fn return_discussion(&self) -> bool {
        to_boolean(self.params().get("return_discussion").unwrap_or("false"))
    }

    fn use_note_serializer(&self) -> bool {
        !this.params().contains_key("html")
    }

    fn params(&self) -> &HashMap<String, String> {
        // Implementation depends on your request parameter handling
        &HashMap::new()
    }
}

impl NotesActions for NotesActionsImpl {
    fn index(&self, req: HttpRequest) -> impl Responder {
        self.set_polling_interval_header(&req);
        self.require_last_fetched_at_header(&req)?;
        self.require_noteable(&req)?;

        let (notes, meta) = self.gather_all_notes();
        let notes = self.prepare_notes_for_rendering(notes);
        let notes = notes
            .into_iter()
            .filter(|n| n.readable_by(&this.current_user))
            .collect::<Vec<_>>();

        let notes = if self.use_note_serializer() {
            self.note_serializer().represent(notes)
        } else {
            notes
                .into_iter()
                .map(|note| self.note_json(&note))
                .collect::<Vec<_>>()
        };

        HttpResponse::Ok().json(meta.merge(serde_json::json!({ "notes": notes })))
    }

    fn create(&self, req: HttpRequest) -> impl Responder {
        self.require_noteable(&req)?;
        self.note_project(&req)?;
        self.check_rate_limit(&req)?;

        let params = web::Json::<CreateNoteParams>::from_request(&req)
            .map_err(|_| HttpResponse::BadRequest())?;

        let note = CreateService::new(&this.project, &this.current_user, &params.0).execute();

        let mut json = serde_json::json!({
            "commands_changes": note.commands_changes()
                .map(|c| c.slice(&["emoji_award", "time_estimate", "spend_time"]))
        });

        if note.is_persisted() && self.return_discussion() {
            json["valid"] = serde_json::json!(true);

            let discussion = note.discussion();
            self.prepare_notes_for_rendering(discussion.notes());
            json["discussion"] = self.discussion_serializer().represent(discussion, this);
        } else {
            self.prepare_notes_for_rendering(vec![note.clone()]);
            json = json.merge(self.note_json(&note));
        }

        let quick_actions = note.quick_actions_status();
        if let Some(quick_actions) = quick_actions {
            json["quick_actions_status"] = serde_json::json!(quick_actions.to_h());
        }

        if note.has_errors() {
            HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": self.errors_on_create(&note)
            }))
        } else if quick_actions.map_or(false, |q| q.error()) {
            HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "quick_actions_status": quick_actions.unwrap().to_h()
            }))
        } else {
            HttpResponse::Ok().json(json)
        }
    }

    fn update(&self, req: HttpRequest) -> impl Responder {
        self.authorize_admin_note(&req)?;

        let note = self.note.as_ref().ok_or_else(|| HttpResponse::NotFound())?;
        let params = web::Json::<UpdateNoteParams>::from_request(&req)
            .map_err(|_| HttpResponse::BadRequest())?;

        let note = UpdateService::new(&this.project, &this.current_user, &params.0).execute(note);

        if note.is_destroyed() {
            return HttpResponse::Gone();
        }

        if note.has_errors() {
            HttpResponse::UnprocessableEntity().json(serde_json::json!({
                "errors": note.errors().join(", ")
            }))
        } else {
            self.prepare_notes_for_rendering(vec![note.clone()]);
            HttpResponse::Ok().json(self.note_json(&note))
        }
    }

    fn destroy(&self, req: HttpRequest) -> impl Responder {
        self.authorize_admin_note(&req)?;

        let note = self.note.as_ref().ok_or_else(|| HttpResponse::NotFound())?;
        if note.editable() {
            DestroyService::new(&this.project, &this.current_user).execute(note);
        }

        HttpResponse::Ok()
    }
}
