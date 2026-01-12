resource "cmccloudv2_dns_zone" "zone_1" {
    domain = "example.com"
    type   = "primary"
}