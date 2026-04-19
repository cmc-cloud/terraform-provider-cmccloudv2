# create mongo backup
resource "cmccloudv2_mongo_backup" "backup1" {
    instance_id   = "b3d915fa-06ae-42a6-847c-cf71fc47a66c"
    name          = "mongo-backup1"
}