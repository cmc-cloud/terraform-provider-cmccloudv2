# create mysql auto backup schedule
resource "cmccloudv2_mysql_autobackup" "autobackup1" {
    instance_id   = "4233c081-0f75-4d9a-83f0-3641a0da1a62"
    schedule_time = "04:10"
    interval      = 1
    max_keep      = 3
}