use crate::controllers::concerns::GroupTree;
use crate::controllers::explore::ExploreApplicationController;
use crate::services::groups::GroupsFinder;
use actix_web::{web, HttpResponse};

const MAX_QUERY_SIZE: i32 = 10_000;

pub struct GroupsController {
    base: ExploreApplicationController,
}

impl GroupsController {
    pub fn new(base: ExploreApplicationController) -> Self {
        Self { base }
    }

    pub async fn index(&self) -> HttpResponse {
        let groups = GroupsFinder::new(self.base.user.clone())
            .execute()
            .await
            .limit(MAX_QUERY_SIZE);

        self.render_group_tree(groups).await
    }
}

impl GroupTree for GroupsController {}
