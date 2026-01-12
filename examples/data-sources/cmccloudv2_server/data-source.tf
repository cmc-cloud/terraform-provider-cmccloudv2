# get server by id
data "cmccloudv2_server" "server1" {  
    server_id = "44fa95ef-5b29-4c81-840f-c2ae310cdbdd"
}

# get server by name
data "cmccloudv2_server" "server1" {  
    name = "ecs-jw73"
}

# get server by ip_address
data "cmccloudv2_server" "server1" {  
    ip_address = "192.168.0.13"
}