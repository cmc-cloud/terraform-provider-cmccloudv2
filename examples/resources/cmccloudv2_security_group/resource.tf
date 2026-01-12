resource "cmccloudv2_security_group" "security_group1" {    
    name 		= "security_group_1" 
    description = "sg rules for public http access"
    stateful 	= true 
}