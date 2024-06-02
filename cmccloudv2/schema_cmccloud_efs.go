package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func efsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_mode": {
			Type:         schema.TypeString,
			Default:      "monthly",
			Optional:     true,
			ValidateFunc: validateBillingMode,
		},
		"capacity": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Capacity in GB",
			ValidateFunc: validation.All(
				validation.IntDivisibleBy(100),
				validation.IntAtLeast(1000),
			),
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			// Default:  "high_performance",
		},
		"protocol_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			// Default:  "nfs",
		},
		"tags": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"endpoint": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"shared_path": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"command_line": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
