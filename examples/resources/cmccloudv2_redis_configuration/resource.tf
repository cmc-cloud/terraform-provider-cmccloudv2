# create redis configuration
resource "cmccloudv2_redis_configuration" "redis_conf_1" {   
    name             = "redis-config-5sac"
    database_engine  = "Redis"
    database_version = "6.0"
    database_mode    = "Master/Slave"
    description      = "default template for master/slave redis database"
    parameters       = {
        zset-max-ziplist-value    = "64"
        zset-max-ziplist-entries  = "128"
        slowlog-log-slower-than   = "10000"
        lua-time-limit            = "5000"
        latency-monitor-threshold = "0"
        hash-max-ziplist-value    = "64"
        hash-max-ziplist-entries  = "512"
        timeout                   = "0"
        slowlog-max-len           = "128"
        notify-keyspace-events    = "Ex"
        set-max-intset-entries    = "512"
        repl-timeout              = "80"
        repl-backlog-size         = "16384"
    }
}