# create a simple volume
resource "cmccloudv2_volume" "volume_1" {  
    name 			= "volume_1" 
    description 	= "volume create from terraform"
    size 			= 20
    type 			= "highio"
    billing_mode 	= "hourly"
    zone 			= "AZ1" 
	tags 			= [{"key": "env", "value": "prod"}]
}  

# create a volume with encryption (this may not working in some zone)
data "cmccloudv2_keymanagement_secret" "key1" {  
	container_id 	= "your_keymanagement_secret_container_id"
    name 			= "key1"
    type 			= "symmetric"
}

resource "cmccloudv2_volume" "volume_2" {  
    name 			= "volume_2" 
    description 	= "volume create from terraform"
    size 			= 20
    type 			= "highio-encryption"
    billing_mode	= "hourly"
    zone 			= "AZ1" 
	secret_id 		= "${data.cmccloudv2_keymanagement_secret.key1.id}"
	tags 			= [{"key": "env", "value": "prod"}]
}