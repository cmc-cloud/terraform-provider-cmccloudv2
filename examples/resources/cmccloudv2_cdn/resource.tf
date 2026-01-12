# create cdn for Origin Host
resource "cmccloudv2_cdn" "cdn1" {
	vod          = false
	name         = "vnexpress site"
	origin_type  = "host"
	domain_or_ip = "vnexpress.net"
	protocol     = "https"
	port         = 443
	origin_path  = "/content/"
}

# create cdn for S3 Host
resource "cmccloudv2_cdn" "cdn2" {
	vod            = false
	name           = "dantri site"
	origin_type    = "s3"
	s3_access_key  = "F7KL3TOAIEKPD4PQ6DRF"
	s3_secret_key  = "J4iwL4PyLPn1y79tg6fQsD2MwZruIEo39k2ckDG0"
	s3_bucket_name = "testbucket"
	s3_region      = "hn-1"
	s3_endpoint    = "s3.hn-1.cloud.cmctelecom.vn"
}

# create cdn for a VOD site 
resource "cmccloudv2_cdn" "cdn3" {
 vod          = true
 name         = "bongda"
 origin_type  = "host"
 domain_or_ip = "bongda.com.vn"
 protocol     = "https"
 port         = 443
}


# create cdn for a VOD site with S3
resource "cmccloudv2_cdn" "cdn4" {
	vod            = true
	name           = "bongda-s3"
	origin_type    = "s3"
	s3_access_key  = "F7KL3TOAIEKPD4PQ6DRF"
	s3_secret_key  = "J4iwL4PyLPn1y79tg6fQsD2MwZruIEo39k2ckDG0"
	s3_bucket_name = "testbucket"
	s3_region      = "hn-1"
	s3_endpoint    = "s3.hn-1.cloud.cmctelecom.vn"
}