# create l7policy with action is REDIRECT_TO_POOL
resource "cmccloudv2_elb_l7policy" "policy1" {
	listener_id  		= "e785611d-a21a-42d3-8062-6ae2cad66b28"
	name         		= "policy-qsav1"  
	action 				= "REDIRECT_TO_POOL"
	position         	= 3
	redirect_pool_id    = "8528ab15-bedf-45bb-843f-abcec42d7e2b"
}

# create l7policy with action is REDIRECT_TO_URL
resource "cmccloudv2_elb_l7policy" "policy2" {
	listener_id  		= "e785611d-a21a-42d3-8062-6ae2cad66b28"
	name         		= "policy-qsav2"  
	action 				= "REDIRECT_TO_URL"
	position         	= 4
	redirect_url  		= "https://google.com"
	redirect_http_code  = 301
	depends_on 			= [cmccloudv2_elb_l7policy.policy1]
}

# create l7policy with action is REJECT
resource "cmccloudv2_elb_l7policy" "policy3" {
	listener_id  		= "e785611d-a21a-42d3-8062-6ae2cad66b28"
	name         		= "policy-qsav3"  
	action 				= "REJECT"
	position         	= 2
	depends_on 			= [cmccloudv2_elb_l7policy.policy2]
}