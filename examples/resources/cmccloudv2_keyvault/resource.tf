data "cmccloudv2_flavor_dbaas" "flavor_dbaas" {
    name = "c5.small.1.db"
}
data "cmccloudv2_volume_type_database" "volume_type_1" {
    description = "High I/O"
}
data "cmccloudv2_flavor_dbaas" "flavor_dbproxy" {
    name = "c5.small.1.db"
}
 
resource "cmccloudv2_keyvault" "keyvault1" {
    name               = "keyvault-2bdr"
    billing_mode       = "monthly"
    version            = "1.21"
    mode               = "ha_cluster"
    zones              = ["AZ1"]
    flavor_id          = "${data.cmccloudv2_flavor_dbaas.flavor_dbaas.id}"
    volume_type        = "${data.cmccloudv2_volume_type_database.volume_type_1.id}"
    volume_size        = 20
    subnet_id          = "fa918710-3310-434a-9dc4-81a44e10bd7f"
    slave_count        = 2
    proxy_flavor_id    = "${data.cmccloudv2_flavor_dbaas.flavor_dbproxy.id}"
    proxy_quantity     = 2
	tags {
        key = "env"
        value = "prod"
    }
}