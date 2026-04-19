# create Master Slave postgres configuration
resource "cmccloudv2_postgres_configuration" "postgres_conf_1" {   
    name             = "postgres-config-terraform1"
    database_version = "15"
    database_mode    = "master_slave" 
    description      = "template for Master Slave postgres database"
    parameters       = {
        "max_stack_depth" = "6200"
    }
}
# create Ha Cluster postgres configuration
resource "cmccloudv2_postgres_configuration" "postgres_conf_2" {   
    name             = "postgres-config-terraform2"
    database_version = "15"
    database_mode    = "ha_cluster" 
    description      = "template for Ha Cluster postgres database"
    parameters       = {
        "wal_writer_flush_after" = "1280"
    }
}

# create standalone postgres configuration
resource "cmccloudv2_postgres_configuration" "postgres_conf_3" {   
    name             = "postgres-config-terraform2"
    database_version = "15"
    database_mode    = "standalone" 
    description      = "template for standalone postgres database"
    parameters       = {
        "autovacuum_vacuum_threshold" = "1000"
    }
}