# create a custom role
resource "cmccloudv2_iam_custom_role" "cr-tf" {
	name        = "cr-tf"
	description = "custom role from terraform"
	content     = file("${path.module}/customrole.json")
}