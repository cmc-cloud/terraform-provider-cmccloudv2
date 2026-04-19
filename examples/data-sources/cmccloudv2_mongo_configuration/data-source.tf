
# get default configuration
data "cmccloudv2_mongo_configuration" "config_default" {
    name       = "default_config_mongodb_6.0_standalone"
    is_default = true
}

# get custom configuration by id
data "cmccloudv2_mongo_configuration" "config_id" {  
    configuration_id = "5f71a7d2-d189-4475-a5cc-e45af49b30b8"
}

# get custom Replica Set configuration by name
data "cmccloudv2_mongo_configuration" "config_replicaset" {
    name          = "mongodb-xoga-config-P8Yuo"
    database_mode = "Replica Set"
}

# get custom Standalone configuration by name
data "cmccloudv2_mongo_configuration" "config_standalone" {
    name          = "mongodb-sa-dev-config-DwVUR"
    database_mode = "Standalone"
}