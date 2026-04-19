# create Replica Set mysql configuration
resource "cmccloudv2_mysql_configuration" "mysql_conf_1" {   
    name             = "mysql-config-terraform1"
    database_version = "15"
    database_mode    = "replica_set" 
    description      = "template for Replica Set mysql database"
    parameters       = {
        "max_write_lock_count" = "18446744073709551615"
        "innodb_ft_min_token_size" = "3"
    }
}

# create Standalone mysql configuration
resource "cmccloudv2_mysql_configuration" "mysql_conf_2" {   
    name             = "mysql-config-terraform2"
    database_version = "15"
    database_mode    = "standalone" 
    description      = "template for Standalone mysql database"
    parameters       = {
        "long_query_time" = "1"
        "sort_buffer_size" = "262144"
    }
}