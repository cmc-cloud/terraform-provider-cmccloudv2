# get custom role by name
data "cmccloudv2_iam_custom_role" "crdev" {
    name = "cr-dev"
}