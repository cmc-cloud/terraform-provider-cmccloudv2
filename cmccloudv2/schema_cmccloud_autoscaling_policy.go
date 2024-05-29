package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func autoscalingHealthCheckPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"interval": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  300,
		},
		"period": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  60,
		},
	}
}
func autoscalingDeletePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"criteria": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"OLDEST_FIRST", "OLDEST_PROFILE_FIRST", "YOUNGEST_FIRST", "RANDOM"}, true),
		},
		"grace_period": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  60,
		},
		"destroy_after_deletion": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},
		"reduce_desired_capacity": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
		"lifecycle_hook_url": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"lifecycle_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Default:  3600,
		},
	}
}
func autoscalingAZPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"zones": {
			Type:     schema.TypeSet,
			Required: true,
			MinItems: 1,
			ForceNew: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}
func autoscalingLoadbalancerPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"lb_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"lb_pool_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"lb_protocol_port": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  443,
			ForceNew: true,
		},
		"as_configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"health_monitor_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
	}
}
func autoscalingScaleInOutPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"scale_number": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"scale_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "One of follow values: CHANGE_IN_CAPACITY,EXACT_CAPACITY,CHANGE_IN_CAPACITY,CHANGE_IN_PERCENTAGE",
			ValidateFunc: validation.StringInSlice([]string{"CHANGE_IN_CAPACITY", "EXACT_CAPACITY", "CHANGE_IN_CAPACITY", "CHANGE_IN_PERCENTAGE"}, true),
		},
		"cooldown": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Length in number of seconds before the actual deletion happens, This param buys an instance some time before deletion",
			ForceNew:    true,
		},
	}
}
