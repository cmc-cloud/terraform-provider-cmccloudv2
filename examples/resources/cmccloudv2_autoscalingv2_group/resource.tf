# create autoscale group do not using Load Balancer
resource "cmccloudv2_autoscalingv2_group" "as_group_1" { 
    name             = "as-group-no-lb"
    zone             = "AZ1"
    min_size         = 1
    max_size         = 1
    desired_capacity = 1
    configuration_id = "9b398fb0-395b-459b-8531-c7ee41e15d30"
	
	scale_up_adjustment_type = "change_in_capacity"
	scale_up_adjustment      = 2
	scale_up_cooldown        = 400
	
	scale_down_adjustment_type = "change_in_capacity"
	scale_down_adjustment      = 1
	scale_down_cooldown        = 400
}  

# create autoscale group using Load Balancer
resource "cmccloudv2_autoscalingv2_group" "as_group_1" { 
	name             = "as-group-vjw2"
	zone             = "AZ1"
	min_size         = 1
	max_size         = 1
	desired_capacity = 1
	configuration_id = "9b398fb0-395b-459b-8531-c7ee41e15d30"
	cooldown         = 600
	
	scale_up_adjustment_type = "change_in_capacity"
	scale_up_adjustment      = 2
	scale_up_cooldown        = 400
	
	scale_down_adjustment_type = "change_in_capacity"
	scale_down_adjustment      = 1
	scale_down_cooldown        = 400
	
	lb_pool_id       = "0e34fbc9-7ab8-4386-a3b9-ab2c1508330f"
	lb_protocol_port = 80
}