# get default configuration
data "cmccloudv2_mysql_configuration" "config_default" {
    name       = "default_config_mysql_8.0_replica_set"
    is_default = true
}

# get custom configuration by id
data "cmccloudv2_mysql_configuration" "config_id" {
    configuration_id = "0ef62cd9-a912-4553-b26c-3f828c70f104"
}

# get custom Replica Set configuration by name
data "cmccloudv2_mysql_configuration" "config_replicaset" {
    name          = "rds_mysql-twbx-config-yqyTN"
    database_mode = "Replica Set"
}

# get custom Standalone configuration by name 
data "cmccloudv2_mysql_configuration" "config_standalone" {
    name          = "rds_mysql-ucnt-config-Oqzng"
    database_mode = "Standalone"
}
