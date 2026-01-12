# create a certificate, two files `waf_cert.cer` & `waf_cert.key` are stored in the current directory
resource "cmccloudv2_waf_cert" "cert1" {
    name        = "cert1"
    cert_name   = "waf_cert.cer"
    cert_data   = file("waf_cert.cer")
    key_name    = "waf_cert.key"
    key_data    = file("waf_cert.key")
    description = "cert for waf"
}