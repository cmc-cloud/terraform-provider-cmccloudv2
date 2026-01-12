resource "cmccloudv2_iam_user" "user" {
	short_name = "dev"
	first_name = "Tuan"
	last_name  = "Nguyen"
	password   = "UGFuq2TqeC"
	email      = "test@example.com"
	enabled    = true
}