# create mysql backup
resource "cmccloudv2_mysql_backup" "backup1" {
    instance_id   = "4233c081-0f75-4d9a-83f0-3641a0da1a62"
    name          = "mysql-backup1"
}
