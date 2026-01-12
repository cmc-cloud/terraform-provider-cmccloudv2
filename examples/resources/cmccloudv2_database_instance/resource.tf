resource "cmccloudv2_database_instance" "database_instance_1" {    
    name              = "db-aie3"
    flavor_id         = "d92cc916-a4cb-4ebf-ae81-801e963465ea"
    zone              = "AZ1"
    source_type       = "new" # backup/instance
    # source_id         = "" # backup_id/instance_id
    datastore_type    = "mysql"
    datastore_version = "8.0.28"
    volume_type       = "highio" # commonio
    volume_size       = 20
    subnets {
	    subnet_id  = "321122f8-8a14-4384-bfec-0306882e9cbf"
		# ip_address = "192.168.4.10"
    }
    enable_public_ip = false
    admin_user       = "rootaa"
    admin_password   = "rOot#12aa"
    billing_mode     = "monthly"
    replicate_count  = 1
    is_public        = true
    # allowed_cidrs    = ""
    # allowed_host     = ""
}