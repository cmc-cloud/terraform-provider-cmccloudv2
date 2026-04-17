resource "cmccloudv2_postgres_database" "postgres_dbprod" {
    instance_id = "e8d3f389-4008-4b54-a1c1-f4152b73878e"
    name        = "dbmain"
    owner       = "admin"
}
