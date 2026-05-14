resource "cmccloudv2_image" "image1" {
    volume_id   = "a4621162-839e-4d06-be8c-93ccd39d7ed7"
    name        = "image-os-ubuntu"
    disk_format = "raw"
    force       = true
}
