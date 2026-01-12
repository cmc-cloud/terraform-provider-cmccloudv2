# create group
resource "cmccloudv2_iam_group" "dev" {
	name        = "dev"
	description = "group from tf"
}