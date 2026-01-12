# set permission `view` only for user `3hr4enp52tvg_dev` on server `464a2fd1-eccb-4568-b20d-993f15355a2a`
resource "cmccloudv2_iam_user_server_permission" "server_read" {
	user_name    = "3hr4enp52tvg_dev"
	server_id    = "464a2fd1-eccb-4568-b20d-993f15355a2a"
	blocked      = false
	allow_view   = true
	allow_edit   = false
	allow_create = false
	allow_delete = false
}