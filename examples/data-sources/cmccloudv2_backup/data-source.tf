# get backup by id
data "cmccloudv2_backup" "backup_1" {  
    backup_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
}

# get backup by name
data "cmccloudv2_backup" "backup_2" {  
    name = "backup-dev"
}

# get latest full backup of a volume
data "cmccloudv2_backup" "backup_3" {  
    volume_id = "2feafa48-8fe7-426f-b35f-aa7761fa97ec"
	status = "available"
	is_latest = true
	is_incremental = false
}