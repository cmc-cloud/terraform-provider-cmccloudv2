package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createAutoscaleConfigurationVolumesElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Volume type, eg: highio/commonio",
			Required:    true,
		},
		"delete_on_termination": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func autoscalingConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"source_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"image", "snapshot", "volume"}, true),
		},
		"source_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"subnet_ids": {
			Type:     schema.TypeList,
			ForceNew: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
				// Required:     true,
				// ValidateFunc: validateUUID,
			},
			Required: true,
			MinItems: 1,
		},
		"use_eip": {
			Type:     schema.TypeBool,
			ForceNew: true,
			Required: true,
		},
		"domestic_bandwidth": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Optional: true,
		},
		"inter_bandwidth": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Optional: true,
		},
		"volumes": {
			Type:     schema.TypeList, // TypeList => (where ordering doesn’t matter), TypeList (where ordering matters).
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: createAutoscaleConfigurationVolumesElementSchema(),
			},
		},
		"security_group_names": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			ForceNew: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"password": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"ecs_group_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
