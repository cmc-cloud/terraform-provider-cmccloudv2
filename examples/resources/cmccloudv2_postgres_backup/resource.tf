# create mysql backup
resource "cmccloudv2_postgres_backup" "backup1" {
    instance_id   = "d4c5ebdb-0447-44d9-b173-375da198a3a7"
    name          = "postgres-backup1"
}