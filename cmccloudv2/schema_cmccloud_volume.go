package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func volumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			// Default:  "",
		},
		"size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validateRegexp(`(?m)^\d{1,10}0$`), // size must be end with 0
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
		},
		"zone": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		// "server_id": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ValidateFunc: validateUUID,
		// },
		"tags": tagSchema(),
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
