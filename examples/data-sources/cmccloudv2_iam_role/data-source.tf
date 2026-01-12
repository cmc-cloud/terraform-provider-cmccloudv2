# get role by name
data "cmccloudv2_iam_role" "ec_full" {
    name = "ec_fullaccess"
}

# get role by title
data "cmccloudv2_iam_role" "ec_full" { 
	title = "EC Administrator" 
}