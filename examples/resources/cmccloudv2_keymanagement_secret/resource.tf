resource "cmccloudv2_keymanagement_secret" "passphrase_1" {
    container_id = "900a7d29-9f06-4eed-ac8f-9bb020ff9ad1"
    name         = "passphrase"
    type         = "passphrase"
    content      = "0123456789aA"
    expiration   = "2026-06-11T16:05:06.277Z"
}