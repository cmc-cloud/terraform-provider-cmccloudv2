data "cmccloudv2_volume_type" "highio" { 
    description = "High I/O" 
}

data "cmccloudv2_volume_type" "ultra" { 
    description = "Ultra High I/O" 
}

resource "cmccloudv2_autoscalingv2_configuration" "as_configuration_1" { 
    name                 		= "as_config_prod"
    source_type          		= "image"
    source_id            		= "23b8854c-581f-40aa-a05a-711d47ac6bd0"
    flavor_id            		= "6aef9929-440a-486a-a487-b1954f56cf87"
    subnet_ids           		= [ "4ca19493-e57f-407f-b3c4-ecb468adea2c" ]
    use_eip              		= true
    domestic_bandwidth   		= 500
    inter_bandwidth      		= 30
    security_group_names 		= [ "default" ]
    password             		= "}2LBgdue5ty"
    ecs_group_id         		= ""
    volumes {
        size                  	= 20
        type                  	= "${data.cmccloudv2_volume_type.highio.name}"
        delete_on_termination 	= true
    }	
    volumes {	
        size                  	= 20
        type                  	= "${data.cmccloudv2_volume_type.ultra.name}"
        delete_on_termination 	= true
    }
}