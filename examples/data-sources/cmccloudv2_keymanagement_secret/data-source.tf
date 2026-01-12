# get keymanagement secret by id
data "cmccloudv2_keymanagement_secret" "redis_config_1" {  
    container_id = "8a26f333-c1b6-4ebf-9ecf-2b155bfaa37d"
}

# get keymanagement secret by name
data "cmccloudv2_keymanagement_secret" "redis_config_1" {  
    name = "secret-a37d"
}