
# get default configuration
data "cmccloudv2_redis_configuration" "config_default" {
    name       = "default_config_redis_7.0_standalone"
    is_default = true
}

# get custom configuration by id
data "cmccloudv2_redis_configuration" "config_id" {  
    configuration_id = "5a910eb8-08f2-4ae5-bfe2-00111de7c5ee"
}

# get custom Cluster configuration by name
data "cmccloudv2_redis_configuration" "config_cluster" {
    name          = "redis-fepx-config-u7Bsu"
    database_mode = "Cluster"
}

# get custom Master/Slave configuration by name
data "cmccloudv2_redis_configuration" "config_master_slave" {
    name          = "redis-jp26-config-7Hymg" 
    database_mode = "Master/Slave"
}

# get custom Standalone configuration by name
data "cmccloudv2_redis_configuration" "config_standalone" {
    name          = "redis-wxyv-config-AYuWS" 
    database_mode = "Standalone"
}