# create a vpc
resource "cmccloudv2_vpc" "vpc_1" {
    name 			= "vpc_1"
    description 	= "VPC HN HC" 
    billing_mode 	= "monthly"
    cidr 			= "192.168.1.0/24"  
	tags 			= [{"key": "env", "value": "prod"}]
}