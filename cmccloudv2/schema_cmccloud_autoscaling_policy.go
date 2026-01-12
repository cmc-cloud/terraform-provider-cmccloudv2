package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func autoscalingHealthCheckPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the autoscaling health check policy",
		},
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Name of action to execute",
			ValidateFunc: validation.StringInSlice([]string{"REBOOT", "REBUILD", "RECREATE"}, false),
		},
		"interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     300,
			Description: "Number of seconds between two adjacent checking",
		},
		"period": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     60,
			Description: "Number of seconds since last node update to wait before checking node health",
		},
	}
}
func autoscalingDeletePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the autoscaling delete policy",
		},
		"criteria": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"OLDEST_FIRST", "OLDEST_PROFILE_FIRST", "YOUNGEST_FIRST", "RANDOM"}, true),
			Description:  "The criteria of the autoscaling delete policy",
		},
		"grace_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     60,
			Description: "The grace period of the autoscaling delete policy",
		},
		"destroy_after_deletion": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
			Description: "The destroy after deletion of the autoscaling delete policy",
		},
		"reduce_desired_capacity": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "The reduce desired capacity of the autoscaling delete policy",
		},
		"lifecycle_hook_url": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The lifecycle hook URL of the autoscaling delete policy",
		},
		"lifecycle_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     3600,
			Description: "The lifecycle timeout of the autoscaling delete policy",
		},
	}
}
func autoscalingAZPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the autoscaling AZ policy",
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
			Description: "The zones of the autoscaling AZ policy",
		},
	}
}
func autoscalingLoadbalancerPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the autoscaling loadbalancer policy",
		},
		"lb_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the autoscaling loadbalancer policy",
		},
		"lb_pool_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the autoscaling loadbalancer pool policy",
		},
		"lb_protocol_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      443,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "The protocol port of the autoscaling loadbalancer policy",
		},
		"as_configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the autoscaling loadbalancer configuration policy",
		},
		"health_monitor_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the autoscaling loadbalancer health monitor policy",
		},
	}
}
func autoscalingScaleInOutPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the autoscaling scale in out policy",
		},
		"scale_number": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The scale number of the autoscaling scale in out policy",
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
