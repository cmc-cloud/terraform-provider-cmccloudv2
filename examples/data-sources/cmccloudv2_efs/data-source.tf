# get efs by id
data "cmccloudv2_efs" "efs1" {  
    efs_id = "07d7eb8f-62ae-4751-968a-3955ef0a6974" 
}

# get efs by name
data "cmccloudv2_efs" "efs1" {  
    name = "efs-88vm" 
}