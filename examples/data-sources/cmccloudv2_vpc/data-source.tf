# get vpc by name
data "cmccloudv2_vpc" "vpc_1" {  
    name = "vpc-a9xj"
}

# get vpc by id
data "cmccloudv2_vpc" "vpc_2" {  
    vpc_id = "8a26f333-c1b6-4ebf-9ecf-2b155bfaa37d"
}

# get vpc by other filter options
data "cmccloudv2_vpc" "vpc_2" {  
    description = "vpc_tf"
    cidr = "192.168.0.0/16"
}