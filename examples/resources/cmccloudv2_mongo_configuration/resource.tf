# create Replica Set mysql configuration
resource "cmccloudv2_mongo_configuration" "mongo_conf_1" {   
    name             = "mongo-config-terraform"
    database_version = "7.0"
    database_mode    = "replica_set" 
    description      = "template for replica set mongo database"
    parameters       = {
        "replication.oplogSizeMB" = "12000"
    }
}

# create Standalone mysql configuration  
resource "cmccloudv2_mongo_configuration" "mongo_conf_2" {   
    name             = "mongo-config-terraform"
    database_version = "7.0"
    database_mode    = "standalone" 
    description      = "template for standalone mongo database"
    parameters       = {
        "security.authorization" = "enabled"
    }
}