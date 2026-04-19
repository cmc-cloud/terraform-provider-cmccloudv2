# create Public elb with flavor "small-lb"
data "cmccloudv2_flavor_elb" "elb_flavor1" {  
    name = "small-lb"
}  
resource "cmccloudv2_elb" "elb_1" {    
    name		 	= "elb_1"
    billing_mode 	= "monthly" 
    zone 			= "AZ1"
    flavor_id 		= "${data.cmccloudv2_flavor_elb.elb_flavor1.id}"
    network_type 	= "public"
    bandwidth_mbps 	= 500 
    description 	= "your description" 
	tags {
        key = "env"
        value = "prod"
    }
}

# create Private elb with flavor "small-lb"
resource "cmccloudv2_elb" "elb_2" {    
    name 			= "elb_2"
    billing_mode 	= "monthly" 
    zone 			= "AZ1"
    flavor_id 		= "${data.cmccloudv2_flavor_elb.elb_flavor1.id}"
    network_type 	= "private"
    subnet_id 		= "d32fa7ba-2a02-4327-80d3-9e17274b9fdd"
    description 	= "your description" 
	tags {
        key = "env"
        value = "prod"
    }
}