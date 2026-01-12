resource "cmccloudv2_kubernetes" "dat_k8s_cluster" {
    name               = "dat_k8s1"
    zone               = "AZ1"
    subnet_id          = "a65cb71b-c7e0-4660-81a6-e46022dbe235"
    docker_volume_size = 20
    docker_volume_type = "highio"
    default_master {
        node_count   = 1
        flavor_id    = "88f0480f-05ac-4e79-accc-167d09e027f0"
        billing_mode = "hourly"
    }
    default_worker {
        node_count     = 1
        max_node_count = 5
        min_node_count = 1
        flavor_id      = "e8de899a-0847-492e-a33f-ef127c97393e"
        billing_mode   = "hourly"
    } 
    labels {
        kube_dashboard_enabled = true
        metrics_server_enabled = true
        npd_enabled            = false
        auto_scaling_enabled   = true
        auto_healing_enabled   = false
        kube_tag               = "v1.26.6-rancher1"
        network_driver         = "calico"
        calico_ipv4pool        = "10.100.0.0/16"
    } 
    create_timeout = 80
    keypair        = "keypair-zn75"
}