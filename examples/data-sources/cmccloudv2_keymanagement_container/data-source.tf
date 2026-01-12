# get keymanagement container by id
data "cmccloudv2_keymanagement_container" "container_1" {  
    container_id = "8a26f333-c1b6-4ebf-9ecf-2b155bfaa37d"
}

# get keymanagement container by name
data "cmccloudv2_keymanagement_container" "container_1" {  
    name = "container-a37d"
}