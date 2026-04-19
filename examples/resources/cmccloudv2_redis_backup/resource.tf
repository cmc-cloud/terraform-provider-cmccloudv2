# create redis backup
resource "cmccloudv2_redis_backup" "backup1" {
    instance_id   = "2dfb8590-73cb-4661-9c18-94e0d83722af"
    name          = "redis-backup1"
}