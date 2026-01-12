resource "cmccloudv2_database_autobackup" "database_autobackup1" {    
    name          = "ab-daily"
    instance_id   = "b7120f46-ea04-4042-959a-0f717c38aba2"
    schedule_time = "04:10"
    interval      = 1
    max_keep      = 3
    incremental   = true
}