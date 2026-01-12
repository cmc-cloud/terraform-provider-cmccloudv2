resource "cmccloudv2_kubernetes_nodegroup" "nodegroup1" {
    cluster_id         = "ee396bfd-1815-4b7a-8570-1141e45b9044"
    name               = "nodegroup1"
    flavor_id          = "e8de899a-0847-492e-a33f-ef127c97393e"
    node_count         = 2
    min_node_count     = 2
    max_node_count     = 3
    billing_mode       = "hourly"
    docker_volume_size = 120
    docker_volume_type = "commonio"
    zone               = "AZ1"
}