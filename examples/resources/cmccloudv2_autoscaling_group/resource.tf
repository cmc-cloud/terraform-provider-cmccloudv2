# create autoscaling group with all policies
resource "cmccloudv2_autoscaling_group" "as_group1" { 
    name                = "as_group1"
    min_size            = 2
    max_size            = 3
    desired_capacity    = 2
    as_configuration_id = "a9497b1d-da0d-4c53-91c3-aa4090efe295"
    policies            = [
        cmccloudv2_autoscaling_health_check_policy.health_check_policy1.id,
        cmccloudv2_autoscaling_az_policy.az_policy1.id,
        cmccloudv2_autoscaling_delete_policy.delete_policy1.id,
        cmccloudv2_autoscaling_lb_policy.lb_policy1.id,
        cmccloudv2_autoscaling_scale_in_policy.scale_in_policy1.id,
        cmccloudv2_autoscaling_scale_out_policy.scale_out_policy1.id
    ]
}