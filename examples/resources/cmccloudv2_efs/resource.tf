# create a simple EFS
resource "cmccloudv2_efs" "efs_1" {
    billing_mode 	= "monthly"
    capacity 		= 1000
    name 			= "efs-aidx"
    subnet_id		= "a7d6c281-64d6-44b6-b02c-cfbf8141adfa"
    type 			= "hdd_standard"
    protocol_type 	= "nfs"
	tags 			= [{"key": "env", "value": "prod"}]
    description 	= "your_description" 
}