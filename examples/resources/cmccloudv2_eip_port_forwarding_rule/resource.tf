resource "cmccloudv2_eip_port_forwarding_rule" "rule1"{
    eip_id              = "9d3e1177-7bc4-4a05-a96e-06a117720965"
    protocol            = "tcp"
    internal_ip_address = "192.168.0.151"
    internal_port_id    = "aa2b0794-3a8e-4c61-a11e-a9dd23af5c04"
    internal_port       = "443"
    external_port       = "443"
    description         = "web port"
}