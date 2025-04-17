use super::renders_member_access::RendersMemberAccess;

pub struct Project {
    pub id: i32,
    pub forks_count: i32,
    pub open_issues_count: i32,
    pub open_merge_requests_count: i32,
    // Add other project fields as needed
}

pub trait RendersProjectsList: RendersMemberAccess {
    fn prepare_projects_for_rendering(&self, projects: &[Project]) -> Vec<Project> {
        let mut prepared_projects = projects.to_vec();

        self.preload_max_member_access_for_collection(&prepared_projects);

        if let Some(_) = self.get_current_user() {
            self.preload_member_roles(&mut prepared_projects);
        }

        // Call the count methods on every project
        for project in &mut prepared_projects {
            project.forks_count();
            project.open_issues_count();
            project.open_merge_requests_count();
        }

        prepared_projects
    }

    fn preload_member_roles(&self, _projects: &mut [Project]) {
        // Overridden in EE version
    }
}

impl Project {
    pub fn forks_count(&mut self) -> i32 {
        self.forks_count
    }

    pub fn open_issues_count(&mut self) -> i32 {
        self.open_issues_count
    }

    pub fn open_merge_requests_count(&mut self) -> i32 {
        self.open_merge_requests_count
    }
}
