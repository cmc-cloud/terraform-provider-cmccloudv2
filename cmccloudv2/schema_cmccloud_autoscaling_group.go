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
			Description:  "Name of the autoscaling group",
		},
		"min_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Minimum number of instances in the autoscaling group",
		},
		"max_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Maximum number of instances in the autoscaling group",
		},
		"desired_capacity": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Desired number of instances in the autoscaling group",
		},
		"as_configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			Description:  "Autoscaling configuration ID",
		},
		"policies": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Description: "List of autoscaling policy IDs",
		},
	}
}
