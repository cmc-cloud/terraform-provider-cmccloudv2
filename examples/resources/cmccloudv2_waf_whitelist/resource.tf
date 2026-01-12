# create a waf rule
resource "cmccloudv2_waf_whitelist" "whitelist_1" {
    waf_id              = "35d40c01-dc76-4478-9864-20f0870586aa"
    message             = "whitelist"
    detection           = "detect str"
    action              = "BLOCK"
    description         = "des"
    match_request_body  = true
    match_get_arguments = true
    match_http_headers  = true
    match_filename      = true
    match_url           = true
    match_name_check    = true
    match_header_var    = "Cookie" # Cookie, Content-Type, User-Agent, Accept-Encoding, Connection
}