resource "cmccloudv2_postgres_autobackup" "autobackup1" {   
    instance_id   = "dfd30170-e27e-4dd6-8647-ba1f9f5debf4"
    schedule_time = "04:10"
    interval      = 1
    max_keep      = 3
}