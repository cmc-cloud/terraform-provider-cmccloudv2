# create LB policy
resource "cmccloudv2_autoscaling_lb_policy" "lb_policy1" {
    name                = "as_group1_lb_policy"
    lb_id               = "f3b07bfa-236e-4f2d-a082-92cbc6c1516b"
    lb_pool_id          = "41fca098-8b10-42ca-8c1c-8cc1622973c1"
    lb_protocol_port    = 443
    as_configuration_id = "a9497b1d-da0d-4c53-91c3-aa4090efe295"
    health_monitor_id   = "871f30c9-c980-4cea-a815-d188be633874"
}