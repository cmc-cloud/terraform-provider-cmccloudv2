# create scale out policy
resource "cmccloudv2_autoscaling_scale_out_policy" "scale_out_policy1" {  
    name         = "as_group1_scale_out_policy"
    scale_number = 1
    scale_type   = "CHANGE_IN_CAPACITY" # CHANGE_IN_CAPACITY, EXACT_CAPACITY, CHANGE_IN_CAPACITY, CHANGE_IN_PERCENTAGE
    cooldown     = 10
}