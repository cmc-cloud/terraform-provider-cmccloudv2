# create a k8s nodegroup without autoscale
data "cmccloudv2_flavor_k8s" "flavor_k8s" {
    name = "c6.small.1.k8s"
}
data "cmccloudv2_volume_type" "ssd" { 
    description = "High I/O (SSD)"
}

resource "cmccloudv2_kubernetesv2_nodegroup" "nodegroup_1" {
	cluster_id         = "ad5a8b1d-0c12-4cd5-b7d9-19a32abfb87f"
	billing_mode       = "monthly"
	name               = "nodegroup_1"
	zone               = "AZ1"
	flavor_id          = "${data.cmccloudv2_flavor_k8s.flavor_k8s.id}"
	key_name           = "keypair-zn75"
	security_group_ids = []
	image_gpu_tag      = "default"
	volume_type        = "${data.cmccloudv2_volume_type.ssd.name}"
	volume_size        = 20
	init_current_node  = 1
	max_pods           = 110
} 

# create a k8s nodegroup with advance options
resource "cmccloudv2_kubernetesv2_nodegroup" "nodegroup_1" {
	cluster_id         = "ad5a8b1d-0c12-4cd5-b7d9-19a32abfb87f"
	billing_mode       = "monthly"
	name               = "nodegroup_1"
	zone               = "AZ1"
	flavor_id          = "${data.cmccloudv2_flavor_k8s.flavor_k8s.id}"
	key_name           = "keypair-zn75"
	security_group_ids = []
	image_gpu_tag      = "default"
	volume_type        = "${data.cmccloudv2_volume_type.ssd.name}"
	volume_size        = 20
	init_current_node  = 1
	max_pods           = 110
	
	#autoscale options
    enable_autoscale = true
    min_node         = 1
    max_node         = 2
	
	#autohealing options
    enable_autohealing           = true
    max_unhealthy_percent        = 80
    node_startup_timeout_minutes = 10
	
	node_metadatas {
		key   = "group"
		value = "cmccloud"
		type  = "label"
	} 
	ntp_servers {
		host     = "vn.pool.ntp.org"
		port     = 123
		protocol = "tcp"
	} 
}