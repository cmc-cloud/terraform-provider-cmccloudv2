# create storage gateway 
resource "cmccloudv2_storage_gateway" "storage_gateway1" {
    name          = "sg-55ys"
    description   = "storage gateway for production"
    protocol_type = "NFS"
    subnet_id     = "7f0dca91-6e37-4d41-820a-e32faec487d3"
    bucket        = "mybucket"
    tags {
        key   = "env"
        value = "prod"
    }
}