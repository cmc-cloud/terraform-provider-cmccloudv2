resource "cmccloudv2_postgres_user" "postgres_useradmin" {
  instance_id = "e8d3f389-4008-4b54-a1c1-f4152b73878e"
  username    = "userdev"
  password    = "5k9DJoK1F"
  permissions = [ "CREATEDB", "CREATEROLE", "LOGIN", "REPLICATION" ]
}