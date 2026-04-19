# create master_slave redis configuration
resource "cmccloudv2_redis_configuration" "redis_conf_1" {   
    name             = "redis-config-5sac"
    database_engine  = "Redis"
    database_version = "6.0"
    database_mode    = "master_slave"
    description      = "template for Master/Slave redis database"
    parameters       = {
        "repl-timeout"    = "64"
        "repl-backlog-size"  = "128"
        "slowlog-log-slower-than"   = "10000"
    }
}
# create cluster redis configuration
resource "cmccloudv2_redis_configuration" "redis_conf_1" {   
    name             = "redis-config-5sac"
    database_engine  = "Redis"
    database_version = "6.0"
    database_mode    = "cluster"
    description      = "template for Cluster redis database"
    parameters       = {
        "set-max-intset-entries"    = "512"
        "repl-backlog-ttl"  = "3600"
        "repl-timeout"   = "60"
    }
}
# create standalone redis configuration
resource "cmccloudv2_redis_configuration" "redis_conf_1" {   
    name             = "redis-config-5sac"
    database_engine  = "Redis"
    database_version = "6.0"
    database_mode    = "standalone"
    description      = "template for Standalone redis database"
    parameters       = {
        "set-max-intset-entries"    = "512"
        "latency-monitor-threshold"  = "0"
        "hash-max-ziplist-entries"   = "512"
    }
}