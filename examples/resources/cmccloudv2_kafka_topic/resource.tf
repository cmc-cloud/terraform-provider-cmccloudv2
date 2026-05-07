resource "cmccloudv2_kafka_topic" "kafka_topic" {
    instance_id        = "c40d86f2-7b91-4059-b04f-8af41ddf60d1"
    name               = "Topic-adex"
    partition_count    = 9
    replication_factor = 2
    rentation_day      = 7
}
