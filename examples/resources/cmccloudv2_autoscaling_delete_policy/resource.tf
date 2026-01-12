# create delete policy
resource "cmccloudv2_autoscaling_delete_policy" "delete_policy1" {  
    name                    = "as_group1_delete_policy"
    criteria                = "OLDEST_FIRST"
    grace_period            = 60
    destroy_after_deletion  = true
    reduce_desired_capacity = true
    lifecycle_hook_url      = ""
    lifecycle_timeout       = 3600
}