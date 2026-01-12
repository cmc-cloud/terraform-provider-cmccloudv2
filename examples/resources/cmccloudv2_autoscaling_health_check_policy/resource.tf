# create health check policy
resource "cmccloudv2_autoscaling_health_check_policy" "health_check_policy1" {  
    name     = "as_group1_health_check1"
    interval = 300
    action   = "RECREATE"
    period   = 100
}