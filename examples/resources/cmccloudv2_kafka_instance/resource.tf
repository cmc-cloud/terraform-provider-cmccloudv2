data "cmccloudv2_flavor_kafka" "flavor_kafka" {
    name = "c6.large.2.kafka"
}
data "cmccloudv2_volume_type" "ssd" {
    description = "Database volume (SSD)"
}
 
# create kafka in Cluster mode
resource "cmccloudv2_kafka_instance" "kafka_instance_cluster" {
    name                = "kafka-2bdr"
    billing_mode        = "monthly"
    version             = "3.7"
    mode                = "cluster"
    broker_quantity     = 3
    zones               = ["AZ1", "AZ3"]
    flavor_id           = data.cmccloudv2_flavor_kafka.flavor_kafka.id
    subnet_id           = "00d793a9-be8c-464d-9163-b27e1058db0b"
    volume_type         = data.cmccloudv2_volume_type_database.ssd.name
    volume_size         = 20
    security_group_ids  = ["2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685"]
    enable_basic_authen = true
    users {
        username = "username-1"
        password = "kBTTPrGOfk"
    }
    users {
        username = "username-2"
        password = "XRvjJWkItG"
    }
}


resource "cmccloudv2_kafka_instance" "kafka_instance_standalone" {
    name                = "kafka-standalone"
    billing_mode        = "monthly"
    version             = "3.7"
    mode                = "single_node"
    zones               = ["AZ1"]
    flavor_id           = data.cmccloudv2_flavor_kafka.flavor_kafka.id
    subnet_id           = "00d793a9-be8c-464d-9163-b27e1058db0b"
    volume_type         = data.cmccloudv2_volume_type_database.ssd.name
    volume_size         = 20
    security_group_ids  = ["2465c8f5-1aa5-4fcd-9ea0-0713d6a1f685"]
    enable_basic_authen = false
}
