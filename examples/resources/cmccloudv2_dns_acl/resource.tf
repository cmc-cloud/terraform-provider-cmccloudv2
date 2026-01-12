# create DNS ACL
resource "cmccloudv2_dns_acl" "acl1" {
	zone_id     = "8dda9b63-3cfa-4b35-98ad-a495e6e1d052"
	record_type = "A"
	domain      = "sub.example.com"
	source_ip   = "1.2.5.0/24"
	action      = "block"  # allow, block
}