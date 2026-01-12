# get DBaAS flavor by name 
data "cmccloudv2_flavor_dbaas" "flavor_ec_1" {
    name = "c6.small.1"  
}

# get DBaAS flavor by id
data "cmccloudv2_flavor_dbaas" "flavor_ec_2" {
    flavor_id = "c9b0f96d-e72d-48f3-b89f-cfb605ab193c" 
}

# get DBaAS flavor by CPU & RAM & Disk
data "cmccloudv2_flavor_dbaas" "flavor_ec_3" {
    cpu = 1
    ram = 2
	disk = 100
}