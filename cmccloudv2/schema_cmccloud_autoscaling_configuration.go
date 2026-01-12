package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createAutoscaleConfigurationVolumesElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Volume size in GB",
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Volume type, eg: highio/commonio",
			Required:    true,
		},
		"delete_on_termination": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "If set to true, delete the volume when the instance is terminated",
		},
	}
}

func autoscalingConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "Name of the autoscaling configuration",
		},
		"source_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"image", "snapshot", "volume"}, true),
			Description:  "Source type of server",
		},
		"source_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "Image/Snapshot/Volume ID, relate to source_type value",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "Flavor ID",
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
			Description: "List of Subnet IDs",
		},
		"use_eip": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Required:    true,
			Description: "Whether to use EIP",
		},
		"domestic_bandwidth": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "Domestic bandwidth, required when use_eip is true",
		},
		"inter_bandwidth": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "Inter bandwidth, required when use_eip is true",
		},
		"volumes": {
			Type:     schema.TypeList, // TypeList => (where ordering doesn’t matter), TypeList (where ordering matters).
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: createAutoscaleConfigurationVolumesElementSchema(),
			},
			Description: "List of volumes",
		},
		"security_group_names": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			ForceNew:    true,
			Description: "List of security group names",
		},
		"key_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Key name",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "User script to run after server created",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: "Password of the server",
		},
		"ecs_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "ECS group ID",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time of the autoscaling configuration",
		},
	}
}
