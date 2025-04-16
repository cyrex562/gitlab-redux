pub mod app;
pub mod diff_file;
pub mod diff_file_header;
pub mod merge_request_diff_file;
mod streaming_error;
pub mod viewers;

pub use app::{AppComponent, DiffFile, DiffsSlice};
pub use diff_file::{DiffFileComponent, DiffFileData};
pub use diff_file_header::{DiffFileHeaderComponent, DiffFileHeaderData};
pub use merge_request_diff_file::{MergeRequest, MergeRequestDiffFileComponent};
pub use streaming_error::StreamingErrorComponent;
pub use viewers::text::LineContentComponent;
