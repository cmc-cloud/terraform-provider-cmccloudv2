# create record with load_balance_type is weighted
resource "cmccloudv2_dns_record" "record1" {
	zone_id           = "08670573-0940-4d9b-9305-caa680f80bbc"
	type              = "A"
	domain            = "example.com"
	ttl               = 300
	load_balance_type = "weighted"
	ips {
        ip     = "1.2.3.4"
        weight = 1
    }
    ips {
        ip     = "5.6.7.8"
        weight = 2
    }
}

# create record with load_balance_type is none
resource "cmccloudv2_dns_record" "record2" {
	zone_id           = "08670573-0940-4d9b-9305-caa680f80bbc"
	type              = "A"
	domain            = "sub.example.com"
	ttl               = 300
	load_balance_type = "none"
	ips {
		ip = "1.2.3.4"
    }
    ips {
		ip = "5.6.7.8"
    }
}