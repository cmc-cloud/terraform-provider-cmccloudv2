
data "cmccloudv2_rds_cluster_flavor" "os_small" {
    name = "db.t1.micro"
}
# create rds_cluster 
resource "cmccloudv2_rds_cluster" "rds_cluster1" {
    billing_mode                  = "monthly"
    name                          = "rds-hrdm"
    flavor_id                     = "${data.cmccloudv2_rds_cluster_flavor.os_small.id}"
    volume_size                   = 20
    db_engine                     = "mysql"
    db_version                    = "8.0.42-33.1"
    subnet_id                     = "7f0dca91-6e37-4d41-820a-e32faec487d3"
    mode                          = "cluster"
    cluster_size                  = 1
    proxy_size                    = 1
    enable_backup                 = false
    enable_pitr                   = true
    backup_schedule               = "0 0 * * *"
    backup_retention              = 3
    enable_storage_autoscaling    = false
    storage_autoscaling_threshold = 70
    storage_autoscaling_increment = 10
}