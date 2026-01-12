# create waf with custom certificate 
resource "cmccloudv2_waf_cert" "cert1" {
    name        = "cert1"
    cert_name   = "waf_cert.cer"
    cert_data   = file("waf_cert.cer")
    key_name    = "waf_cert.key"
    key_data    = file("waf_cert.key")
    description = "cert for waf"
}
resource "cmccloudv2_waf" "waf1" { 
    domain                 = "google.com"
    mode                   = "DETECT"
    real_server            = "116.104.138.10"
    protocol               = "HTTPS"
    port                   = 443
    certificate_id         = cmccloudv2_waf_cert.cert1.id
    description            = "desc"
    load_balance_enable    = false
    load_balance_method    = "ip_hash"
    load_balance_keepalive = 1000
}