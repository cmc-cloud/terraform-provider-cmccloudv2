# get eip by id
data "cmccloudv2_eip" "elb1" { 
    eip_id = "af262476-df43-4ad6-a626-37b77e23ed14" 
}

# get eip by ip address
data "cmccloudv2_eip" "elb2" { 
    ip_address = "203.171.29.45" 
}