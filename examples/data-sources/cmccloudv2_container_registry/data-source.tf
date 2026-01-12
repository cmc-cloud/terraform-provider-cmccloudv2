# get container registry by id
data "cmccloudv2_container_registry" "ecs_group_1" {  
    devops_project_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
	container_registry_id = "491ff792-9627-409b-9eb2-e444582fa45a"
}

# get container registry by name
data "cmccloudv2_container_registry" "ecs_group_1" {  
    devops_project_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
	name = "cr-kx37"
}