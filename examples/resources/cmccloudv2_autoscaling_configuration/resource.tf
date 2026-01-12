resource "cmccloudv2_autoscaling_configuration" "as_configuration1" { 
    name               = "as_config_terraform"
    source_type        = "image"
    source_id          = "c9cd9428-84a1-4e77-946a-0f8c44a2eccc"
    flavor_id          = "f195840e-1fe1-4816-af3b-54d7390a00a2"
    subnet_ids         = [ "52257938-65ae-422b-97f4-211db1de00a1" ]
    use_eip            = true
    domestic_bandwidth = 500
    inter_bandwidth    = 10
    volumes {
        size                  = 20
        type                  = "highio"
        delete_on_termination = true
    }
    volumes {
        size                  = 100
        type                  = "commonio"
        delete_on_termination = true
    }
    security_group_names = [ "sg-cv8g", "default" ]
    # key_name             = ""
    # user_data            = ""
    password             = "}2LBgdue5ty"
    ecs_group_id         = "bfdcd02a-1ffe-4e24-9cc5-09a0a6689923"
}