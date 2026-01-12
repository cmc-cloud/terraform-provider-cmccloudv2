# get redis configuration by id
data "cmccloudv2_redis_configuration" "redis_config_1" {  
    configuration_id = "8a26f333-c1b6-4ebf-9ecf-2b155bfaa37d"
}

# get redis configuration by name
data "cmccloudv2_redis_configuration" "redis_config_2" {  
    name = "config-m97t"
}

# get redis configuration by mode
data "cmccloudv2_redis_configuration" "redis_config_3" {  
    database_mode = "Standalone"
}