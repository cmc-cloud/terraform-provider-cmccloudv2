# get port by ip address + server
data "cmccloudv2_port" "port1" {
    server_id  = "06f07c77-8c5e-4faa-8532-b526684dfe61"
    ip_address = "10.10.1.180"
}

# get port by id
data "cmccloudv2_port" "port2" {
    port_id = "8a26f333-c1b6-4ebf-9ecf-2b155bfaa37d"
}

# get first private port of server
data "cmccloudv2_port" "port3" {
    server_id = "06f07c77-8c5e-4faa-8532-b526684dfe61"
    is_public = false
}
