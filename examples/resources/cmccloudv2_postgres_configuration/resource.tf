resource "cmccloudv2_postgres_configuration" "postgres_conf_1" {   
    name             = "postgres-config-5sac"
    database_version = "17"
    database_mode    = "standalone" 
    description      = "default template for standalone postgres database"
    parameters       = {
        "join_collapse_limit"= "8"
        "enable_incremental_sort"= "on"
        "wal_log_hints"= "off"
        "logical_decoding_work_mem"= "65536"
        "wal_keep_size"= "4096"
        "hot_standby_feedback"= "off"
		"max_parallel_workers_per_gather" = "2"
		"synchronous_commit" = "on"
    }
}