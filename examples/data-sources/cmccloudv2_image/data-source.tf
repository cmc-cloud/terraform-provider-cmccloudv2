# get image by name
data "cmccloudv2_image" "window2016" {
    visibility = "public" # public, shared, private
    name = "Windows Server 2016" 
}

# get image by id
data "cmccloudv2_image" "window2016" {
    image_id = ""
}