package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func autoscalingGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"min_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"max_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"desired_capacity": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"as_configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
		},
		"policies": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
		},
	}
}
