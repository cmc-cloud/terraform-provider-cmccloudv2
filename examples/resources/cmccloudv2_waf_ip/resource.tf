# create a waf IP
resource "cmccloudv2_waf_ip" "ip1" {
    waf_id      = "4f50b772-97a7-4177-ab18-2b84629de4d3"
    ip          = "116.104.138.50"
    type        = "deny" # deny, ignore
    description = "desc"
}