resource "cmccloudv2_database_configuration" "database_config_terraform" {    
    name              = "db-config-terraform"
    description       = "Database configuration created from terraform"
    datastore_type    = "mysql"
    datastore_version = "8.0.28"
    parameters {
        key   = "autocommit"
        value = "1"
    }
    parameters {
        key   = "auto_increment_increment"
        value = "1"
    }
    parameters {
        key   = "character_set_server"
        value = "utf8"
    }
}