# create `WEB_APPLICATION_TESTS` Vulnerability Assessment Scan and run immediately
resource "cmccloudv2_va" "va_1" {
    name        = "va-ke89"
    type        = "WEB_APPLICATION_TESTS"
    target      = "https://google.com"
    description = "scan vulnerability"
}

# create `BASIC_NETWORK_SCAN` Vulnerability Assessment Scan and run at specific time
resource "cmccloudv2_va" "va_1" {
    name        = "va-dip3"
    type        = "BASIC_NETWORK_SCAN"
    schedule    = "2026-01-25 01:00:50"
    target      = "https://google.com"
    description = "run at night"
}