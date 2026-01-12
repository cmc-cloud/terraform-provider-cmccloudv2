# Create a rule allows outbound IPv4 UDP traffic on port 80 from instances that belong to your_security_group_id to any destination on the Internet.
resource "cmccloudv2_security_group_rule" "rule1" {
	security_group_id 	= "your_security_group_id"
	direction 			= "egress" 
	protocol 			= "udp"
	port_range_min 		= 80
	port_range_max 		= 80
	cidr 				= "0.0.0.0/0" 
	ether_type 			= "IPv4"
}

# Create a rule allows incoming IPv4 TCP traffic on port 80 (HTTP) to instances that belong to your_security_group_id, from any IPv4 address on the Internet.
resource "cmccloudv2_security_group_rule" "rule2" {
	security_group_id 	= "your_security_group_id"
	direction 			= "ingress" 
	protocol 			= "tcp"
	port_range_min 		= 80
	port_range_max		= 80
	cidr				= "0.0.0.0/0" 
	ether_type 			= "IPv4"
}

# Create a rule allows incoming TCP traffic to servers that belong to `your_security_group_id`, on ports from 10000 to 65535, only if the traffic originates from servers that are members of another security group (another_security_group_id), and only for IPv4 traffic.
resource "cmccloudv2_security_group_rule" "rule3" {
	security_group_id 	= "your_security_group_id"
	direction 			= "ingress"
	protocol 			= "tcp" 
	port_range_min 		= 10000
	port_range_max 		= 65535
	remote_group_id 	= "another_security_group_id"
	ether_type 			= "IPv4"
}