# create a simple interface, attach to a server, a random IP address will be assigned
resource "cmccloudv2_server_interface" "interface_1" {
	server_id 	= "90df1634-798d-4ba2-acb1-ae2a17a9cb94"
	subnet_id 	= "5d6159b5-8903-4b6c-b841-eeed9594b4ce"
}

# create a simple interface with specific ip address, this ip address must belong to the subnet
resource "cmccloudv2_server_interface" "interface_2" {
	server_id 	= "90df1634-798d-4ba2-acb1-ae2a17a9cb94"
	subnet_id 	= "5d6159b5-8903-4b6c-b841-eeed9594b4ce"
	ip_address 	= "192.168.1.10"
}