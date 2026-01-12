# get ecs group by id
data "cmccloudv2_ecs_group" "ecs_group_1" {  
    ecs_group_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
}

# get ecs group by name
data "cmccloudv2_ecs_group" "ecs_group_1" {  
    name = "ecs-group-hfy8"
}