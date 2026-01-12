# get elb by id
data "cmccloudv2_elb" "elb1" { 
    eip_id = "af262476-df43-4ad6-a626-37b77e23ed14" 
}

# get elb by name
data "cmccloudv2_elb" "elb2" { 
    name = "elb-k4d5" 
}

# get elb by ip_address
data "cmccloudv2_elb" "elb3" { 
    ip_address = "203.171.29.45" 
}