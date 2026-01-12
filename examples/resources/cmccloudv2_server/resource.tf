# Create instance with SSH key
resource "cmccloudv2_server" "example_01" {	
	name           	= "example-01"
    billing_mode 	= "monthly" 
    zone 			= "AZ1"
    flavor_id 		= "2b9dea6e-f628-446b-91f9-1122cff5a418" 
    source_id 		= "a9cb8f5e-b32a-4e89-b2d4-6dc2a25d0918"
	source_type 	= "image" 
	volume_name		= "root-example-01"
    volume_size 	= 30
    volume_type 	= "ssd"  
    subnet_id 		= "be679ec3-e661-4010-8e9e-f8c041a68dc7"
	key_name 		= "keypair-rm74" 
}  

# Create instance with password
resource "cmccloudv2_server" "example_02" {	
	name            = "example-02"
    billing_mode 	= "monthly" 
    zone 			= "AZ1"
    flavor_id 		= "2b9dea6e-f628-446b-91f9-1122cff5a418" 
    source_id 		= "a9cb8f5e-b32a-4e89-b2d4-6dc2a25d0918"
	source_type 	= "image" 
	volume_name		= "root-2ds6"
    volume_size 	= 30
    volume_type 	= "ssd"  
    subnet_id 		= "be679ec3-e661-4010-8e9e-f8c041a68dc7"
	password 		= "UGFuq1TqeC@"
}  
 
# Create instance from backup/snapshot
resource "cmccloudv2_server" "example_03" {	
	name            = "example-03"
    billing_mode 	= "monthly" 
    zone 			= "AZ1"
    flavor_id 		= "2b9dea6e-f628-446b-91f9-1122cff5a418" 
    source_id 		= "58e458bc-c04e-4bc7-bb12-1b4dbc115a7c" # backup id
	source_type 	= "backup" 
	volume_name		= "root-example-03"
    volume_size 	= 30
    volume_type 	= "ssd"  
    subnet_id 		= "be679ec3-e661-4010-8e9e-f8c041a68dc7"
	password 		= "UGFuq1TqeC@"
}  


# Create instance with advance optiopns: ip address & server group & security group & start script (user_data)
resource "cmccloudv2_server" "example_04" {	
	name              		= "example-04"
    billing_mode 			= "monthly" 
    zone 					= "AZ1"
    flavor_id 				= "2b9dea6e-f628-446b-91f9-1122cff5a418" 
    source_id 				= "a9cb8f5e-b32a-4e89-b2d4-6dc2a25d0918"
	source_type 			= "image" 
	volume_name				= "root-example-04"
    volume_size 			= 30
    volume_type 			= "ssd"  
    subnet_id 				= "be679ec3-e661-4010-8e9e-f8c041a68dc7"
	password 				= "UGFuq1TqeC@"
	ip_address 				= "192.168.0.55" 
    ecs_group_id 			= "bfdcd02a-1ffe-4e24-9cc5-09a0a6689923"
    user_data 				= ""
	delete_on_termination 	= false
	tags 					= [{"key": "env", "value": "prod"}]
    security_group_names 	= [ "sg-cv8g", "default" ]
}