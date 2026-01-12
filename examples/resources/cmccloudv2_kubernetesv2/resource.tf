# create k8s cluster without addons
data "cmccloudv2_flavor_k8s" "master_flavor" {
    name = "c6.large.2.k8s"
}
resource "cmccloudv2_kubernetesv2" "k8s_1" {
    name               = "k8s_1"
    zone               = "AZ1"
    subnet_id          = "ac67cf25-205d-440a-82de-1017f72653b7"
    kubernetes_version = "v1.28.9"
    master_flavor_name = "c6.large.2.k8s"
    master_count       = 1
    cidr_block_pod     = "10.100.0.0/16"
    cidr_block_service = "10.254.0.0/16"
    network_driver     = "calico"
} 

# create k8s cluster with addons
resource "cmccloudv2_kubernetesv2" "k8s_2" {
    name               = "k8s_2"
    zone               = "AZ1"
    subnet_id          = "ac67cf25-205d-440a-82de-1017f72653b7"
    kubernetes_version = "v1.28.9"
    master_flavor_name = "c6.large.2.k8s"
    master_count       = 1
    cidr_block_pod     = "10.100.0.0/16"
    cidr_block_service = "10.254.0.0/16"
    network_driver     = "calico"
	
    enable_autohealing   = true
    enable_monitoring    = true
    enable_autoscale     = true
    autoscale_max_node   = 50
    autoscale_max_ram_gb = 500
    autoscale_max_core   = 500
} 


# create k8s cluster with advance options
resource "cmccloudv2_kubernetesv2" "k8s_3" {
    name               = "k8s_3"
    zone               = "AZ1"
    subnet_id          = "ac67cf25-205d-440a-82de-1017f72653b7"
    kubernetes_version = "v1.28.9"
    master_flavor_name = "c6.large.2.k8s"
    master_count       = 1
    cidr_block_pod     = "10.100.0.0/16"
    cidr_block_service = "10.254.0.0/16"
    network_driver     = "calico"
	
	network_driver_mode = "native-routing"
	ntp_servers {
		host     = "vn.pool.ntp.org"
		port     = 123
		protocol = "tcp"
	} 
}