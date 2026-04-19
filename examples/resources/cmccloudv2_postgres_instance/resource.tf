data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c6.large.2.db"
}
data "cmccloudv2_volume_type_database" "volume_type_1" {
    description = "Database volume (SSD)"
}
data "cmccloudv2_flavor_dbproxy" "flavor_dbproxy" {
    name = "c6.small.dbproxy"
}

## create postgres in Master/Slave mode
data "cmccloudv2_postgres_configuration" "config_master_slave" {
    name          = "rds_postgres-8cqb-config-3aJVv"
	is_default 	  = false
    database_mode = "Master Slave"
}
resource "cmccloudv2_postgres_instance" "postgres_instance_masterslave" {
    name               = "postgres-2bdr"
    billing_mode       = "monthly"
    version            = "15"
    mode               = "master_slave"
    zones              = ["AZ1", "AZ3"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type        = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size        = 20
    subnet_id          = "f9de9530-1000-43f1-98ec-2c97dc49020f"
    admin_password     = "5k9DJoK1FqDQS"
    #backup_id          = ""    
    configuration_id   = "${data.cmccloudv2_postgres_configuration.config_master_slave.id}" # configuration template must have the same mode + database version
    port               = 5432
    retention_period   = 3
    slave_count        = 1 
	tags {
        key = "env"
        value = "prod"
    }
}

## create postgres in Ha Cluster mode
data "cmccloudv2_postgres_configuration" "config_ha_cluster" {
    name          = "rds_postgres-xzdq-config-hZA4R"
	is_default 	  = false
    database_mode = "Ha Cluster"
}
resource "cmccloudv2_postgres_instance" "postgres_instance_ha_cluster" {
    name               = "postgres-2bdr"
    billing_mode       = "monthly"
    version            = "15"
    mode               = "ha_cluster"
    zones              = ["AZ1", "AZ3"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type        = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size        = 20
    subnet_id          = "f9de9530-1000-43f1-98ec-2c97dc49020f"
    admin_password     = "5k9DJoK1FqDQS" 
    configuration_id   = "${data.cmccloudv2_postgres_configuration.config_ha_cluster.id}" # configuration template must have the same mode + database version
    port               = 5432
    retention_period   = 3
    slave_count        = 1
    proxy_quantity     = 2
    proxy_flavor_id    = "${data.cmccloudv2_flavor_dbproxy.flavor_dbproxy.id}"
    tags {
        key = "env"
        value = "prod"
    }
}

## create postgres in Standalone mode
data "cmccloudv2_postgres_configuration" "config_ha_standalone" {
    name          = "rds_postgres-dzvf-config-E8mz5"
	is_default 	  = false
    database_mode = "Standalone"
}
resource "cmccloudv2_postgres_instance" "postgres_instance_standalone" {
    name               = "postgres-2bdr"
    billing_mode       = "monthly"
    version            = "15"
    mode               = "standalone"
    zones              = ["AZ1"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type        = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size        = 20
    subnet_id          = "f9de9530-1000-43f1-98ec-2c97dc49020f"
    admin_password     = "5k9DJoK1FqDQS"
    #backup_id          = ""
    configuration_id   = "${data.cmccloudv2_postgres_configuration.config_ha_standalone.id}" # configuration template must have the same mode + database version
    port               = 5432
	retention_period   = 3
    tags {
        key = "env"
        value = "prod"
    }
}