data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c6.large.2.db"
}
data "cmccloudv2_volume_type_database" "volume_type_1" {
    description = "Database volume (SSD)"
}
# create mysql in Replica Set mode
data "cmccloudv2_mysql_configuration" "config_replicaset" {
    name          = "rds_mysql-twbx-config-yqyTN"
    database_mode = "Replica Set"
}
resource "cmccloudv2_mysql_instance" "mysql_instance_replica1" {
    name                  = "mysql-2bdr"
    billing_mode          = "monthly"
    version               = "8.0"
    mode                  = "replica_set"
    zones                 = ["AZ1", "AZ3"]
    flavor_id             = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type           = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size           = 20
    subnet_id             = "f9de9530-1000-43f1-98ec-2c97dc49020f"
    configuration_id      = "${data.cmccloudv2_mysql_configuration.config_replicaset.id}"
    #backup_id             = ""
	tags {
		key = "env"
		value = "prod"
	}
	tags {
		key = "ver"
		value = "v1"
	}
    quantity_of_secondary = 2
}

# create mysql in Standalone mode
data "cmccloudv2_mysql_configuration" "config_standalone" {
    name          = "rds_mysql-ucnt-config-Oqzng"
    database_mode = "Standalone"
}
resource "cmccloudv2_mysql_instance" "mysql_instance_standalone" {
    name                  = "mysql-2bdr"
    billing_mode          = "monthly"
    version               = "8.0"
    mode                  = "standalone"
    zones                 = ["AZ1"]
    flavor_id             = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type           = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size           = 20
    subnet_id             = "f9de9530-1000-43f1-98ec-2c97dc49020f"
    configuration_id      = "${data.cmccloudv2_mysql_configuration.config_standalone.id}"
	tags {
		key = "env"
		value = "prod"
	}
}