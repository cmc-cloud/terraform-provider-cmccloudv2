# get subnet by name
data "cmccloudv2_subnet" "subnet1" {  
    name = "subnet-ekw6"
}

# get subnet by id
data "cmccloudv2_subnet" "subnet2" {  
    subnet_id = "321122f8-8a14-4384-bfec-0306882e9cbf"
}

# get subnet by using other filter options
data "cmccloudv2_subnet" "subnet3" {   
    vpc_id = "53c71330-21d8-47f3-8df9-49d7d345b873"
    cidr = "192.168.4.0/24"
    gateway_ip = "192.168.4.1"
}