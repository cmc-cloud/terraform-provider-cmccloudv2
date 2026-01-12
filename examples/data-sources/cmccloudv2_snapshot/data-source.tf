# get snapshot by id
data "cmccloudv2_snapshot" "snapshot_1" {  
    snapshot_id = "3a27e174-2ad7-4988-bfd4-c985e8abf0b4"
}

# get snapshot by name
data "cmccloudv2_snapshot" "snapshot_1" {  
    name = "snapshot-dex4"
}

# get latest snapshot of volume
data "cmccloudv2_snapshot" "snapshot-latest" {
    volume_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974"
    status = "available"
    is_latest = true 
}