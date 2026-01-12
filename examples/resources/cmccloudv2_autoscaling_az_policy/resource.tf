# create AZ policy
resource "cmccloudv2_autoscaling_az_policy" "az_policy1" {  
    name  = "as_group1_az_policy"
    zones = [ "AZ1" ]
}