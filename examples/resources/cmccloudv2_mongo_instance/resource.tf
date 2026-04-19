data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c6.large.2.db"
}
data "cmccloudv2_volume_type_database" "volume_type_1" {
    description = "Database volume (SSD)"
} 
# create mongo in Replica Set mode
data "cmccloudv2_mongo_configuration" "config_replicaset" {
    name          = "mongodb-xoga-config-P8Yuo"
    database_mode = "Replica Set"
}
resource "cmccloudv2_mongo_instance" "mongo_instance_replica1" {
    name                  = "mongo-2bdr"
    billing_mode          = "monthly"
    version               = "6.0"
    mode                  = "replica_set"
    zones                 = ["AZ1", "AZ2"]
    flavor_id             = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type           = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size           = 20
    subnet_id             = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    configuration_id      = "${data.cmccloudv2_mongo_configuration.config_replicaset.id}"
    security_group_ids    = ["2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685"]
    #backup_id             = ""
    quantity_of_secondary = 2
	tags {
        key = "env"
        value = "prod"
    }
}

# create mongo in Standalone mode
data "cmccloudv2_mongo_configuration" "config_standalone" {
    name          = "mongodb-sa-dev-config-DwVUR"
    database_mode = "Standalone"
}
resource "cmccloudv2_mongo_instance" "mongo_instance_standalone" {
    name                  = "mongo-2bdr"
    billing_mode          = "monthly"
    version               = "6.0"
    mode                  = "standalone"
    zones                 = ["AZ1"]
    flavor_id             = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type           = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size           = 20
    subnet_id             = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    configuration_id      = "${data.cmccloudv2_mongo_configuration.config_standalone.id}"
    security_group_ids    = ["2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685"]
	tags {
        key = "env"
        value = "prod"
    }
}