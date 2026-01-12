# get ELB flavor by name 
data "cmccloudv2_elb_flavor" "flavor_elb_1" {
    flavor_id = "af262476-df43-4ad6-a626-37b77e23ed14"
}

# get ELB flavor by id
data "cmccloudv2_elb_flavor" "flavor_elb_2" {
    name = "small-lb" 
}