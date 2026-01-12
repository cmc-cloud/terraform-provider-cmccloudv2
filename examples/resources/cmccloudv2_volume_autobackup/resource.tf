# Create auto backup run at 04:30 every day and keep 7 backups in chains
resource "cmccloudv2_server" "autobackup_01" {	
	name			= "auto-backup-1"
	volume_id		= "642109d2-49ab-48b1-8fc7-8f1909c33a1c"
	schedule_time	= "04:30"
	incremental		= false
	interval		= 1
	max_keep		= 7 
}