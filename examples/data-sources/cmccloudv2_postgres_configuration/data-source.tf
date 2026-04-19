# get default configuration
data "cmccloudv2_postgres_configuration" "config_default" {
    name       = "default_config_postgresql_17_standalone"
    is_default = true
}

# get custom configuration by id
data "cmccloudv2_postgres_configuration" "config_id" {  
    configuration_id = "5f71a7d2-d189-4475-a5cc-e45af49b30b8"
}

# get custom Ha Cluster configuration by name
data "cmccloudv2_postgres_configuration" "config_hacluster" {
    name          = "rds_postgres-xzdq-config-hZA4R"
    database_mode = "Ha Cluster"
}

# get custom Master Slave configuration by name
data "cmccloudv2_postgres_configuration" "config_master_slave" {
    name          = "rds_postgres-r4rv-config-K0qVg" 
    database_mode = "Master Slave"
}

# get custom Standalone configuration by name
data "cmccloudv2_postgres_configuration" "config_ha_standalone" {
    name          = "rds_postgres-dzvf-config-E8mz5" 
    database_mode = "Standalone"
}