# attach eip with a port
resource "cmccloudv2_eip_port" "eip_port" {    
    eip_id 			= "f036e873-4579-4575-8ebc-2364a7b0e522" 
    port_id 		= "8de5d47e-ff2b-4390-8333-fbf3d775cd1a"
    fix_ip_address 	= "192.168.0.10" 
}