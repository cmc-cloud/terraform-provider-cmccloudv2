data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c6.large.2.db"
}
data "cmccloudv2_volume_type_database" "volume_type_1" {
    description = "Database volume (SSD)"
} 
data "cmccloudv2_mongo_configuration" "config_master_slave" {
    name          = "mongo-hbxm-config-XqCAA"
    database_mode = "master/slave"
}
# create mongo in Replica Set mode
resource "cmccloudv2_mongo_instance" "mongo_instance_replica1" {
    name                  = "mongo-2bdr"
    billing_mode          = "monthly"
    version               = "6.0"
    mode                  = "replica_set"
    zones                 = ["AZ1", "AZ2"]
    flavor_id             = data.cmccloudv2_flavor_dbaas.flavor_dbaas.id
    volume_type           = data.cmccloudv2_volume_type_database.volume_type_1.id
    volume_size           = 20
    subnet_id             = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    configuration_id      = data.cmccloudv2_mongo_configuration.config_master_slave.id
    security_group_ids    = ["2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685"]
    backup_id             = ""
    quantity_of_secondary = 2
}

# create mongo in Standalone mode
resource "cmccloudv2_mongo_instance" "mongo_instance_masterslave" {
    name                  = "mongo-2bdr"
    billing_mode          = "monthly"
    version               = "6.0"
    mode                  = "standalone"
    zones                 = ["AZ1"]
    flavor_id             = data.cmccloudv2_flavor_dbaas.flavor_dbaas.id
    volume_type           = data.cmccloudv2_volume_type_database.volume_type_1.id
    volume_size           = 20
    subnet_id             = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    configuration_id      = data.cmccloudv2_mongo_configuration.config_master_slave.id
    security_group_ids    = ["2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685"] 
}