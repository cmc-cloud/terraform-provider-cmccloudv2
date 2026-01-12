data "cmccloudv2_flavor_dbaas" "flavor_dbaas" { 
    name = "c6.large.2.db"
} 
data "cmccloudv2_redis_configuration" "config_master_slave" {   
    name          = "redis-hbxm-config-XqCAA"
    database_mode = "master/slave"
} 
# create redis in Cluster mode
resource "cmccloudv2_redis_instance" "redis_instance_masterslave" {
    name               = "redis-2bdr"
    billing_mode       = "monthly"
    database_engine    = "Redis"
    database_version   = "6.0"
    zones              = ["AZ3", "AZ2"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_size        = 20
    subnet_id          = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    password           = "5k9DJoK1FqDQShwgocmqkZ23DIlFOaa"
    security_group_ids = [ "2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685" ]
	
	database_mode          = "Cluster"
	replicas               = 1
	redis_configuration_id = "${data.cmccloudv2_redis_configuration.config_master_slave.id}"
}

# create redis in Master/Slave mode
resource "cmccloudv2_redis_instance" "redis_instance_masterslave" {
    name               = "redis-2bdr"
    billing_mode       = "monthly"
    database_engine    = "Redis"
    database_version   = "6.0"
    zones              = ["AZ3", "AZ2"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_size        = 20
    subnet_id          = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    password           = "5k9DJoK1FqDQShwgocmqkZ23DIlFOaa"
    security_group_ids = [ "2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685" ]
	
    database_mode = "Master/Slave"
}

# create redis in Standalone mode & using backup
resource "cmccloudv2_redis_instance" "redis_instance_masterslave" {
    name               = "redis_terraform3"
    billing_mode       = "monthly"
    database_engine    = "Redis"
    database_version   = "6.0"
    zones              = ["AZ1"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_size        = 20
    subnet_id          = "036f3b55-2ff8-4350-9fc0-4baf12deca03"
    password           = "5k9DJoK1FqDQShwgocmqkZ23DIlFOaa"
    security_group_ids = [ "2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685" ]
	
	database_mode = "Standalone"
	backup_id     = "162efde6-caf4-4e5f-b0aa-24fa27c17bfc"
}