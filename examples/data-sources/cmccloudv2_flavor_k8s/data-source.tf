# get k8s flavor by name 
data "cmccloudv2_flavor_k8s" "flavor_k8s_1" {
    name = "c6.small.1"  
}

# get k8s flavor by id
data "cmccloudv2_flavor_k8s" "flavor_k8s_2" {
    flavor_id = "c9b0f96d-e72d-48f3-b89f-cfb605ab193c" 
}

# get k8s flavor by CPU & RAM & Disk
data "cmccloudv2_flavor_k8s" "flavor_k8s_3" {
    cpu = 1
    ram = 2
	disk = 50
}