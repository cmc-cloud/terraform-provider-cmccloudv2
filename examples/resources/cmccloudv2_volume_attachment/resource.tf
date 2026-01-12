resource "cmccloudv2_volume_attachment" "volume_attachment1" {
    volume_id 				= "b116eaac-129c-4e45-9c66-f3e1297b5f5b"
    server_id 				= "85ba4449-f4ad-4272-878b-51e93c218321"
    delete_on_termination 	= false
}