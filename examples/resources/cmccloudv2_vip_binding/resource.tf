# bind vip to a port by id
resource "cmccloudv2_vip_binding" "vip_bind_port1" {
    vip_id = "8de5d47e-ff2b-4390-8333-fbf3d775341d"
    port_id = "ed34625a-5adc-4abb-b112-17e631423ca8"
} 

# bind vip to a port of new server
data "cmccloudv2_volume_type" "ssd" {
    description = "High I/O (SSD)"
}
resource "cmccloudv2_server" "example_01" {
    name                 = "example-01"
    billing_mode         = "monthly"
    zone                 = "AZ1"
    flavor_id            = "2b9dea6e-f628-446b-91f9-1122cff5a418"
    source_id            = "a9cb8f5e-b32a-4e89-b2d4-6dc2a25d0918"
    source_type          = "image"
    volume_name          = "root-example-01"
    volume_size          = 30
    volume_type          = "${data.cmccloudv2_volume_type.ssd.id}"
    subnet_id            = "be679ec3-e661-4010-8e9e-f8c041a68dc7"
    key_name             = "keypair-rm74"
    security_group_names = ["default"]
}
resource "cmccloudv2_vip_binding" "vip_bind_vm1" {
    vip_id = "8de5d47e-ff2b-4390-8333-fbf3d775341d"
    port_id = "${cmccloudv2_server.example_01.interface_id}"
}

# bind vip to a port of existing server by id
# get first private port of server
data "cmccloudv2_port" "port2" {
    server_id = "06f07c77-8c5e-4faa-8532-b526684dfe61"
    is_public = false
}
resource "cmccloudv2_vip_binding" "vip_bind_vm1" {
    vip_id = "8de5d47e-ff2b-4390-8333-fbf3d775341d"
    port_id = "${cmccloudv2_port.port2.id}"
}
