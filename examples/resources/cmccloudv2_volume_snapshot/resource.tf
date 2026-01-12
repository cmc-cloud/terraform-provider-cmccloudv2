resource "cmccloudv2_volume_snapshot" "snapshot_1" {
    name 		= "snapshot-1"
    volume_id 	= "642109d2-49ab-48b1-8fc7-8f1909c33a1c"
    force 		= true
}