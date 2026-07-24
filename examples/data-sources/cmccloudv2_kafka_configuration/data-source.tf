# get default configuration
data "cmccloudv2_kafka_configuration" "config_default" {
    name       = "default_config_kafka_3.7_single_node"
    is_default = true
}

# get custom configuration by id
data "cmccloudv2_kafka_configuration" "config_id" {  
    configuration_id = "5f71a7d2-d189-4475-a5cc-e45af49b30ae"
}

# get custom Cluster configuration by name
data "cmccloudv2_kafka_configuration" "config_cluster" {
    name          = "kafka-xoga-config-ie9d"
    database_mode = "Cluster"
}

# get custom Single-node configuration by name
data "cmccloudv2_kafka_configuration" "config_single_node" {
    name          = "kafka-sa-dev-config-ac43"
    database_mode = "Single-node"
}