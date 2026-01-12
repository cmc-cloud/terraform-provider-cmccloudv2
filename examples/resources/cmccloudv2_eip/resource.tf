#create a simple EIP
resource "cmccloudv2_eip" "eip_terraform" {    
    billing_mode 		= "monthly" 
    domestic_bandwidth 	= 500 
    inter_bandwidth 	= 10
    description 		= "eip created from terraform" 
}

#create a EIP with advance options
resource "cmccloudv2_eip" "eip_terraform" {    
    billing_mode 		= "monthly" 
    domestic_bandwidth 	= 500 
    inter_bandwidth 	= 10
    description 		= "eip created from terraform"
	tags 	 			= [{"key": "env", "value": "prod"}]
    dns_domain 			= "example.com."
    dns_name 			= "my-ip" 
}