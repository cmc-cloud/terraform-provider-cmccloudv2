# get port
data "cmccloudv2_port" "test" {
    server_id = "06f07c77-8c5e-4faa-8532-b526684dfe61"
    is_public = false
}
data "cmccloudv2_security_group" "default" {
    name = "default"
}
# config port
resource "cmccloudv2_port_config" "port1" {
    port_id               = data.cmccloudv2_port.test.id
    # only pass the attributes you want to update
    name                  = "test_port"
    port_security_enabled = true
    security_group_ids    = ["${data.cmccloudv2_security_group.default.id}"]
    allowed_address_pairs {
        ip_address  = "1.3.4.5"
        mac_address = "fa:16:3e:00:e4:70"
    }
}