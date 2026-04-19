# create mysql database
resource "cmccloudv2_mysql_database" "db1" {
    instance_id   = "4233c081-0f75-4d9a-83f0-3641a0da1a62"
    name          = "db1"
    character_set = "utf8mb4"
    collation     = "utf8mb4_unicode_ci"
}
