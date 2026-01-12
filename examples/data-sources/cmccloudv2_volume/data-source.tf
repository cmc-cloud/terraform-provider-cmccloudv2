# get volume by name
data "cmccloudv2_volume" "volume_1" { 
    name = "ev-k8s3"
}

# get root volume of server
data "cmccloudv2_volume" "volume_root" { 
    server_id = "44fa95ef-5b29-4c81-840f-c2ae310cdbdd"
    bootable = true
}

# get volume by id
data "cmccloudv2_volume" "volume_root" { 
    volume_id = "5d0fecd3-728d-46b1-98bd-93d1cefba685"
}