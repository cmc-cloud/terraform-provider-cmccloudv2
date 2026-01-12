package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createAutoscaleV2ConfigurationVolumesElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The size of the volume",
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the volume, eg: highio/commonio",
			Required:    true,
			ForceNew:    true,
		},
		"delete_on_termination": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "If true, the volume will be deleted when the instance is terminated. Default is true",
		},
	}
}

func autoscalingV2ConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the autoscaling configuration",
		},
		"source_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"image", "snapshot", "volume"}, true),
			Description:  "Type of source: image, snapshot or volume",
		},
		"source_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The source ID of image/snapshot/volume, depends on source_type",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The flavor ID of the desired flavor for the server.",
		},
		"subnet_ids": {
			Type:     schema.TypeList,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				// Required:     true,
				// ValidateFunc: validateUUID,
			},
			Required:    true,
			MinItems:    1,
			Description: "The ID of subnets to associate with the server.",
		},
		"use_eip": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Required:    true,
			Description: "If true, the EIP will be used after the server is created",
		},
		"domestic_bandwidth": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "Domestic bandwidth in Mbps, required if use_eip is true",
		},
		"inter_bandwidth": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "International bandwidth in Mbps, required if use_eip is true",
		},
		"volumes": {
			Type:     schema.TypeList, // TypeList => (where ordering doesn’t matter), TypeList (where ordering matters).
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: createAutoscaleV2ConfigurationVolumesElementSchema(),
			},
			Description: "The volumes for the server",
		},
		"security_group_names": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			ForceNew:    true,
			Description: "An array of one or more security group names to associate with the server",
		},
		"key_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The name of a key pair to put on the server",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "User data of server, it will be executed when the server is created",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: "Password of server, it will be used to login to the server.",
		},
		"ecs_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The ECS group ID of the autoscaling configuration",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the autoscaling configuration",
		},
	}
}
