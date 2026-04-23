data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c6.large.2.db"
}
data "cmccloudv2_volume_type_database" "volume_type_1" {
    description = "High I/O (SSD)"
}

# create AZ gateway api
resource "cmccloudv2_apigateway" "azapi_gateway1" {
    name          = "azapi_gateway_1"
    mode          = "standalone"
    zones         = ["AZ1"]
    flavor_id     = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    subnet_id     = "f9de9530-1000-43f1-98ec-2c97dc49020f"
    volume_type   = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size   = 20
    public_access = true
    bandwidth     = 500
    tags {
        key   = "env"
        value = "prod"
    }
}
