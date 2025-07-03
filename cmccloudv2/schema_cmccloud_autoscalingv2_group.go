package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func autoscalingGroupV2Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
		},
		"zone": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
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
		"configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"lb_pool_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
		},
		"lb_protocol_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validatePortNumber,
			RequiredWith: []string{"lb_pool_id"},
		},
		"cooldown": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      600,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "Cooldown time in seconds between scaling actions.",
		},
		"scale_up_adjustment_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "change_in_capacity",
			ValidateFunc: validation.StringInSlice([]string{"change_in_capacity", "percent_change_in_capacity"}, false),
			Description:  "Type of scaling adjustment for scale-up.",
		},
		"scale_up_adjustment": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Value to adjust for scale-up.",
		},
		"scale_up_cooldown": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      300,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "Cooldown period after scale-up in seconds.",
		},
		"scale_down_adjustment_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "change_in_capacity",
			ValidateFunc: validation.StringInSlice([]string{"change_in_capacity", "percent_change_in_capacity"}, false),
			Description:  "Type of scaling adjustment for scale-down.",
		},
		"scale_down_adjustment": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "Value to adjust for scale-down.",
		},
		"scale_down_cooldown": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      300,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "Cooldown period after scale-down in seconds.",
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status_reason": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
