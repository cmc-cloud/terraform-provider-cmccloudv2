# Listener with Layer 4 Protocol TCP/UDP/SCTP
resource "cmccloudv2_elb_listener" "listener_1"{
    elb_id 					= "f054ad90-1d14-461b-958c-983d8a16c6a2"
    name 					= "listener_1"
    description 			= "your_listener_description"
    protocol 				= "TCP"
    protocol_port 			= "80"
    default_pool_id 		= "fcd14e84-578f-46b7-ba2e-b88862903332"
	
    timeout_member_connect 	= 5000
    timeout_member_data 	= 50000
    connection_limit 		= -1
    allowed_cidrs 			= [] 
	
    # two options bellow not avaiable when protocol is UDP/SCTP
    timeout_client_data 	= 50000 
    timeout_tcp_inspect 	= 0
}

# create listener with Layer 7 Protocol - HTTP
resource "cmccloudv2_elb_listener" "listener_2"{
    elb_id 					= "f054ad90-1d14-461b-958c-983d8a16c6a2"
    name 					= "listener_2"
    description 			= "your_listener_description"
    protocol 				= "HTTP"
    protocol_port 			= "80"
    default_pool_id 		= "fcd14e84-578f-46b7-ba2e-b88862903332"
    timeout_member_connect 	= 5000
    timeout_member_data 	= 50000
    connection_limit 		= -1
    allowed_cidrs 			= [] 
    timeout_client_data 	= 50000 
    timeout_tcp_inspect 	= 0
	
	# options for layer 7 protocol
    x_forwarded_for 		= false
    x_forwarded_port 		= true
    x_forwarded_proto 		= true
}

# create listener with Layer 7 Protocol - TERMINATED_HTTPS
data "cmccloudv2_certificate" "cert1" {    
    name = "certificate-xstm" 
}

resource "cmccloudv2_elb_listener" "listener_3"{
    elb_id 						= "f054ad90-1d14-461b-958c-983d8a16c6a2"
    name 						= "listener_3"
    description 				= "your_listener_description"
    protocol 					= "HTTP"
    protocol_port 				= "80"
    default_pool_id 			= "fcd14e84-578f-46b7-ba2e-b88862903332"
    timeout_member_connect 		= 5000
    timeout_member_data 		= 50000
    connection_limit 			= -1
    allowed_cidrs 				= [] 
    timeout_client_data 		= 50000 
    timeout_tcp_inspect 		= 0 
	
    x_forwarded_for 			= false
    x_forwarded_port 			= true
    x_forwarded_proto 			= true
		
    default_tls_container_ref  	= "${data.cmccloudv2_certificate.cert1.secret_ref}"
}