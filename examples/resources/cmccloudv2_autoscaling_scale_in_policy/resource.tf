# create scale in policy
resource "cmccloudv2_autoscaling_scale_in_policy" "scale_in_policy1" {  
    name         = "as_group1_scale_in_policy"
    scale_number = 1
    scale_type   = "CHANGE_IN_CAPACITY"         // CHANGE_IN_CAPACITY,EXACT_CAPACITY,CHANGE_IN_CAPACITY,CHANGE_IN_PERCENTAGE
    cooldown     = 100
}