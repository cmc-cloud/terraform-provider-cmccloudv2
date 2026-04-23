data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c6.large.2.db"
}
data "cmccloudv2_volume_type" "volume_type_1" {
    description = "High I/O (SSD)"
}

# create AZ gateway api
resource "cmccloudv2_apigateway" "azapi_gateway1" {
    name          = "azapi_gateway2"
    mode          = "standalone"
    zones         = ["AZ1"]
    flavor_id     = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    subnet_id     = "1264314a-162e-4165-abf1-fcb0212bc1f5"
    volume_type   = "${data.cmccloudv2_volume_type.volume_type_1.id}"
    volume_size   = 20
    public_access = true
    bandwidth     = 500
    tags {
        key   = "env"
        value = "prod"
    }
}
