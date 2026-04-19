# create redis auto backup schedule
resource "cmccloudv2_redis_autobackup" "autobackup1" {
    instance_id   = "2dfb8590-73cb-4661-9c18-94e0d83722af"
    schedule_time = "04:10"
    interval      = 1
    max_keep      = 3
}