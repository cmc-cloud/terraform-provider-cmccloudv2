# create mysql user with alter and create permissions on all databases & tables
resource "cmccloudv2_mysql_user" "user1" {
	instance_id = "4233c081-0f75-4d9a-83f0-3641a0da1a62"
	username  = "user1"
	password = "AaioeuroieEd343"
	hosts = [ "%" ]
	user_permissions {
		database = "*"
		table = "*"
		permissions = [ "alter", "create" ]
	}
}