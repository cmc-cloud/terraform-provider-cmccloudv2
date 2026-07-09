# create vip with manual ip address
resource "cmccloudv2_vip" "vip1" {
    subnet_id  = "8de5d47e-ff2b-4390-8333-fbf3d775341a"
    ip_address = "192.168.0.10"
}
# create vip with auto ip address
resource "cmccloudv2_vip" "vip2" {
    subnet_id  = "8de5d47e-ff2b-4390-8333-fbf3d775341a" 
}
