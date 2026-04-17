resource "cmccloudv2_rds_cluster_user" "user1" {
    cluster_id = "e996c7c1-9464-4813-80e3-54a9397e8365"
    name       = "user1"
    host       = "localhost"
    password   = "TqcESmLu2d60%"
    databases  = ["dbprod1"]
}