# create Health Monitor using PING/TCP/TLS-HELLO/UDP-CONNECT/SCTP protocol
resource "cmccloudv2_elb_healthmonitor" "health_monitor_1"{
	name    			= "health_monitor_1"
	pool_id 			= "fcd14e84-578f-46b7-ba2e-b88862903332"
	type    			= "PING"
	
    max_retries_down 	= 5
    delay            	= 5
    max_retries      	= 5
    timeout          	= 5
}

# create Health Monitor using HTTP/HTTPS protocol
resource "cmccloudv2_elb_healthmonitor" "health_monitor_2"{
	name    			= "health_monitor_2"
	pool_id 			= "fcd14e84-578f-46b7-ba2e-b88862903332"
	type    			= "HTTP"
	
    max_retries_down 	= 5
    delay            	= 5
    max_retries      	= 5
    timeout          	= 5
   
    http_method    		= "GET"
    expected_codes 		= "200-209"
    url_path       		= "/"
}