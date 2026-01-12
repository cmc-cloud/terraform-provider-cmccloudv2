# create scale up trigger. If the minimum CPU usage of instances in the Auto Scaling Group is greater than 80% for 2 consecutive checks, evaluated every 5 minutes, then scale up
resource "cmccloudv2_autoscalingv2_scale_trigger" "as_group1_trigger_up" { 
	group_id    = "f41c10b0-427d-44fc-ab31-89db464e7697"
	name        = "as-trigger-as_group1-scale-up"
	function    = "min"
	metric      = "cpu"
	comparator  = ">"
	threadhold  = 80
	interval    = 5
	times       = 2
	action      = "scale_up"
	enabled     = true
	description = "desc"
}

# create scale down trigger. If the minimum CPU usage of the Auto Scaling Group is less than or equal to 50%, for 2 consecutive checks, evaluated every 5 minutes, then scale down 
resource "cmccloudv2_autoscalingv2_scale_trigger" "as_group1_trigger_down" { 
	group_id    = "f41c10b0-427d-44fc-ab31-89db464e7697"
	name        = "as-trigger-as_group1-scale-down"
	function    = "min"
	metric      = "cpu"
	comparator  = "<="
	threadhold  = 50
	interval    = 5
	times       = 2
	action      = "scale_down"
	enabled     = true
	description = "desc"
}