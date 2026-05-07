resource "cmccloudv2_mongo_database" "mongodb" {
    instance_id = "4233c081-0f75-4d9a-83f0-3641a0da1a62"
    name        = "mongodb"
    collections = ["coll1", "coll2"]
}
