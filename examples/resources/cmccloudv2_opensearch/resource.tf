data "cmccloudv2_opensearch_flavor" "os_small" { 
    name = "os.small.1"
} 
data "cmccloudv2_opensearch_dashboard_flavor" "dash_small" { 
    name = "dash.small"
}  
# create opensearch 
resource "cmccloudv2_opensearch" "opensearch1" { 
    billing_mode                  = "monthly"
    name                          = "os-55ys"   
    version                       = "2.19.5"
    flavor_id                     = "${data.cmccloudv2_opensearch_flavor.os_small.id}"
    dashboard_flavor_id           = "${data.cmccloudv2_opensearch_dashboard_flavor.dash_small.id}"
    volume_size                   = 20   
    admin_password                = "SV<9wsb7rhhbapb" 
    node_count                    = 2
    enable_isolate_master         = true
    master_count                  = 3
    dashboard_replicas            = 1
    enable_snapshot               = true
    snapshot_creation_cron        = "0 2 * * *"
    snapshot_timezone             = "Asia/Ho_Chi_Minh"
    rentation_max_age             = 14
    rentation_min_count           = 3
    rentation_max_count           = 10
    lb_subnet_id                  = "7f0dca91-6e37-4d41-820a-e32faec487d3"
    enable_storage_autoscaling    = true
    storage_autoscaling_threshold = 80
    storage_autoscaling_increment = 10
    storage_autoscaling_max       = 5000
    tags {
        key = "env"
        value = "prod"
    }
}