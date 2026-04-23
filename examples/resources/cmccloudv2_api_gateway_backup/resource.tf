# create AZ gateway api backup
resource "cmccloudv2_api_gateway_backup" "backup1" {
    name        = "backup_1"
    instance_id = "c47cfe8d-65d4-4b55-b169-adac69f5bdad"
}
