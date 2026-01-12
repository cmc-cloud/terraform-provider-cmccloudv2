# get group by name
data "cmccloudv2_iam_group" "dev" {
    name = "dev"
}