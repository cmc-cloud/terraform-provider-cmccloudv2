# create simple ELB pool
resource "cmccloudv2_elb_pool" "pool_1"{
    elb_id 				= "40b37c15-3888-4bf4-8fb5-396d2363beac"
    name 				= "pool_1"
    description 		= "my pool 1" 
    protocol 			= "HTTP"
    algorithm 			= "ROUND_ROBIN"
    session_persistence = "NONE"
}

# create simple ELB pool with session_persistence is APP_COOKIE
resource "cmccloudv2_elb_pool" "pool_2"{
    elb_id 				= "40b37c15-3888-4bf4-8fb5-396d2363beac"
    name				= "pool_2"
    description 		= "my pool 2" 
    protocol 			= "HTTP"
    algorithm 			= "ROUND_ROBIN"
    session_persistence = "APP_COOKIE"
    cookie_name 		= "your_cookie_name" 
}

# create simple ELB pool with tls support
resource "cmccloudv2_elb_pool" "pool_3"{
    elb_id 				= "40b37c15-3888-4bf4-8fb5-396d2363beac"
    name				= "pool_3"
    description 		= "my pool 3" 
    protocol 			= "HTTP"
    algorithm 			= "ROUND_ROBIN"
    session_persistence = "NONE"
    tls_enabled 		= true
    tls_ciphers 		= "ciphers text"
    tls_versions 		= [ "TLSv1.2", "TLSv1.3" ] 
}