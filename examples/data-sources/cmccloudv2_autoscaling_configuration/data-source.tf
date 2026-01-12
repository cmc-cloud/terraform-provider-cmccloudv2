# get configuration by id
data "cmccloudv2_autoscaling_configuration" "confg_1" {  
    autoscaling_configuration_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
}

# get configuration by name
data "cmccloudv2_autoscaling_configuration" "backup_2" {  
    name = "config-vst6"
}