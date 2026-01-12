# get user by username
data "cmccloudv2_iam_user" "dev" {
    username = "3hr4enp52tvg_dev"
}

# get user by email
data "cmccloudv2_iam_user" "dev" {
	email = "test@example.com"
}