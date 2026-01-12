# get EC flavor by name 
data "cmccloudv2_flavor_ec" "flavor_ec_1" {
    name = "c6.small.1"  
}

# get EC flavor by id
data "cmccloudv2_flavor_ec" "flavor_ec_2" {
    flavor_id = "c9b0f96d-e72d-48f3-b89f-cfb605ab193c" 
}

# get EC flavor by CPU & RAM
data "cmccloudv2_flavor_ec" "flavor_ec_3" {
    cpu = 1
    ram = 2
}