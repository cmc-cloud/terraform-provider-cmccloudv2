# attach member to a pool
resource "cmccloudv2_elb_pool_member" "member_1"{
    pool_id			= "b9d75072-c2f1-41d5-900b-ac2bbe55ec65"
    name 			= "ecs-aunx"
    address 		= "192.168.0.4"
    protocol_port 	= 80
    weight			= 1
    monitor_address = "192.168.0.4"
    monitor_port 	= 80
}