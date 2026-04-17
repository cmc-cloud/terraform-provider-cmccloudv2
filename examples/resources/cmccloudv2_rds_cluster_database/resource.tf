
resource "cmccloudv2_rds_cluster_database" "database1" {
    cluster_id    = "095c0455-23ed-455c-b477-1be7515a31c2"
    name          = "dbprod"
    character_set = "utf8mb4"
    collation     = "utf8mb4_unicode_ci"
}
