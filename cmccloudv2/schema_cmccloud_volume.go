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
			Description:  "Name for the volume. Changing this updates the volume's name.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the volume",
		},
		"size": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validateRegexp(`(?m)^\d{1,10}0$`), // size must be end with 0
			Description:  "The size of the volume in GB, block on 10GB/100GB base on type of volume",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The type of the volume",
		},
		"secret_id": {
			Type:         schema.TypeString,
			ValidateFunc: validateUUID,
			Optional:     true,
			ForceNew:     true,
			Description:  "The secret ID of the volume, this is used to encrypt the volume",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
		},
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The availability zone for the volume. Changing this creates a new volume.",
		},
		"backup_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The backup/snapshot ID from which to create the volume. Changing this creates a new volume",
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
