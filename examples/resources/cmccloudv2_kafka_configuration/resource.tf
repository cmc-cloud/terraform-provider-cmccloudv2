# create master_slave kafka configuration
resource "cmccloudv2_kafka_configuration" "kafka_conf_1" {   
    name             = "kafka-config-5sac"
    database_version = "3.7"
    database_mode    = "single_node"
    description      = "template for Single Node kafka database"
    parameters       = {
        "transaction.state.log.segment.bytes"    = "104857600"
        "transactional.id.expiration.ms"   = "604800000"
    }
}
# create cluster kafka configuration
resource "cmccloudv2_kafka_configuration" "kafka_conf_2" {   
    name             = "kafka-config-a38d"
    database_version = "3.7"
    database_mode    = "cluster"
    description      = "template for Cluster kafka database"
    parameters       = {
        "replica.lag.time.max.ms"    = "30000"
        "log.index.interval.bytes"   = "4096"
    }
} 