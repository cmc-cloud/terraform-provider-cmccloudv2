# Create a rule allows outbound IPv4 UDP traffic on port 80 from instances that belong to your_security_group_id to any destination on the Internet.
resource "cmccloudv2_security_group_rule" "rule1" {
	security_group_id 	= "7815b87b-1973-4e9f-9ae6-e38a858d90f5"
	direction 			= "egress" 
	protocol 			= "udp"
	port_range_min 		= 80
	port_range_max 		= 80
	cidr 				= "0.0.0.0/0" 
	ether_type 			= "IPv4"
}

# Create a rule allows incoming IPv4 TCP traffic on port 80 (HTTP) to instances that belong to your_security_group_id, from any IPv4 address on the Internet.
resource "cmccloudv2_security_group_rule" "rule2" {
	security_group_id 	= "7815b87b-1973-4e9f-9ae6-e38a858d90f5"
	direction 			= "ingress" 
	protocol 			= "tcp"
	port_range_min 		= 80
	port_range_max		= 80
	cidr				= "0.0.0.0/0" 
	ether_type 			= "IPv4"
}

# Create a rule allows incoming TCP traffic to servers that belong to `your_security_group_id`, on ports from 10000 to 65535, only if the traffic originates from servers that are members of another security group (another_security_group_id), and only for IPv4 traffic.
resource "cmccloudv2_security_group_rule" "rule3" {
	security_group_id 	= "7815b87b-1973-4e9f-9ae6-e38a858d90f5"
	direction 			= "ingress"
	protocol 			= "tcp" 
	port_range_min 		= 10000
	port_range_max 		= 65535
	remote_group_id 	= "ac540e27-f442-4196-a91f-4c46a7494106"
	ether_type 			= "IPv4"
}