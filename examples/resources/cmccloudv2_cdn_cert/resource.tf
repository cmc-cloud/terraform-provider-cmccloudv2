# create a cdn certificate, two files `cdn_cert.cer` & `cdn_cert.key` are stored in the current directory
resource "cmccloudv2_cdn_cert" "cert" {
    cert_name = "cdn_cert.cer"
    cert_data = file("cdn_cert.cer")
    key_name  = "cdn_cert.key"
    key_data  = file("cdn_cert.key")
}