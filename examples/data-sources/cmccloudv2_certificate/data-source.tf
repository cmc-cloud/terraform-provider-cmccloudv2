# get container registry by id
data "cmccloudv2_certificate" "cert_1" {  
    certificate_id = "fe416a11-c03a-4a34-8874-632ba0324941"
}

# get container registry by name
data "cmccloudv2_certificate" "cert_2" {  
    name = "cert-a49t"
}