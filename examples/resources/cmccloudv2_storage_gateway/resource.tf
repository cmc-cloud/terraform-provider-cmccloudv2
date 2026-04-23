# create storage gateway 
resource "cmccloudv2_storage_gateway" "storage_gateway1" {
   name          = "sg-55ys1"
   description   = "storage1 gateway for production"
   protocol_type = "nfs"
   subnet_id     = "0899fb16-f32a-416f-8618-ebe8a5af725d"
   bucket        = "mybucket1"
   tags {
       key   = "env"
       value = "test"
   }
}