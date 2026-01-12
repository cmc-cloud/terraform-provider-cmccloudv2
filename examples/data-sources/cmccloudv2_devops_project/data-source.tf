# get devops project by id
data "cmccloudv2_devops_project" "project_1" {  
    devops_project_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
}

# get devops project by name
data "cmccloudv2_devops_project" "project_1" {  
    name = "project-dev"
}