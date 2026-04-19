# create mongo auto backup schedule
resource "cmccloudv2_mongo_autobackup" "autobackup1" {
    instance_id   = "b3d915fa-06ae-42a6-847c-cf71fc47a66c"
    schedule_time = "04:10"
    interval      = 1
    max_keep      = 3
}