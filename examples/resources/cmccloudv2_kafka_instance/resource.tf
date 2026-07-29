data "cmccloudv2_flavor_kafka" "flavor_kafka" {
    name = "c6.large.2.kafka"
}

data "cmccloudv2_volume_type_database" "ssd" {
    description = "Database volume (SSD)"
}

data "cmccloudv2_kafka_configuration" "config_cluster" {
    name          = "kafka-xoga-config-ie9d"
    database_mode = "Cluster"
}

# get custom Single-node configuration by name
data "cmccloudv2_kafka_configuration" "config_single_node" {
    name          = "kafka-sa-dev-config-ac43"
    database_mode = "Single-node"
}
# create kafka in Cluster mode
resource "cmccloudv2_kafka_instance" "kafka_instance_cluster" {
    name                = "kafka-2bdr"
    billing_mode        = "monthly"
    version             = "3.7"
    mode                = "cluster"
    broker_quantity     = 3
    zones               = ["AZ1", "AZ3"]
    flavor_id           = "${data.cmccloudv2_flavor_kafka.flavor_kafka.id}"
    subnet_id           = "00d793a9-be8c-464d-9163-b27e1058db0b"
    volume_type         = "${data.cmccloudv2_volume_type_database.ssd.id}"
    volume_size         = 20
    enable_basic_authen = true
    security_group_ids  = ["00d793a9-be8c-464d-9163-b27e1058db0b"]
    configuration_id    = "${data.cmccloudv2_kafka_configuration.config_cluster.id}"
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
    flavor_id           = "${data.cmccloudv2_flavor_kafka.flavor_kafka.id}"
    subnet_id           = "00d793a9-be8c-464d-9163-b27e1058db0b"
    volume_type         = "${data.cmccloudv2_volume_type_database.ssd.id}"
    volume_size         = 20
    configuration_id    = "${data.cmccloudv2_kafka_configuration.config_single_node.id}"
    enable_basic_authen = false
}
