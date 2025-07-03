package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func autoscalingGroupV2ScaleTriggerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"scale_up", "scale_down"}, false),
			Description:  "The scaling action to perform. Allowed values: scale_up, scale_down.",
		},
		"function": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"min", "max", "sum", "count", "avg", "last"}, false),
			Description:  "The aggregation function to use. Allowed values: min, max, sum, count, avg, last.",
		},
		"metric": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"cpu", "memory", "net_in", "net_out", "io_read", "io_write", "io_read_ops", "io_write_ops"}, false),
			Description:  "The metric to monitor. Allowed values: cpu, memory, net_in, net_out, io_read, io_write, io_read_ops, io_write_ops.",
		},
		"comparator": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{">", ">=", "<", "<="}, false),
			Description:  "Comparator for threshold. Allowed values: >, >=, <, <=",
		},
		"threadhold": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "Threshold value for the metric. Must be >= 0.",
		},
		"interval": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Interval in seconds to evaluate the metric. Must be >= 1.",
		},
		"times": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Number of evaluation periods. Must be >= 1.",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether the scale trigger is enabled.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the scale trigger.",
		},
	}
}
