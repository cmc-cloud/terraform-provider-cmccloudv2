data "cmccloudv2_devops_project" "project_default"{
    name = "project-6k4"
    type = "owner" // owner, share
}

resource "cmccloudv2_container_registry_repo" "repo_1" {
    devops_project_id = "${data.cmccloudv2_devops_project.project_default.id}"
    name = "repository-1"
}