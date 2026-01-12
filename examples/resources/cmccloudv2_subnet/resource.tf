resource "cmccloudv2_subnet" "subnet_1" {
    name 			= "subnet_1"
    vpc_id 			= "b116eaac-129c-4e45-9c66-f3e1297b5f5b"
    enable_dhcp 	= true
	ip_version 		= 4
    gateway_ip 		= "192.168.1.1"
    cidr 			= "192.168.1.0/24"
	tags 			= [{"key": "env", "value": "prod"}]
    dns_nameservers = [ "183.91.10.1", "183.91.10.2" ]
    allocation_pools {
        start 	= "192.168.1.2"
        end 	= "192.168.1.254"
    } 
    host_routes {
        destination = "192.168.1.0/24"
        nexthop 	= "192.168.1.1"
    }
    host_routes {
        destination = "192.168.2.0/24"
        nexthop 	= "192.168.2.1"
    }
}