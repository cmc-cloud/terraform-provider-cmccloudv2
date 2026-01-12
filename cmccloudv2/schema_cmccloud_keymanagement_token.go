package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func keymanagementtokenSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_ids": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Description: "The IDs of the key management containers",
		},

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The description of the key management token",
		},

		"expiration": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
			Description:  "The expiration time of the key management token, example: `2026-06-11T16:05:06.277Z`",
		},

		"token": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The token of the key management token",
		},

		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the key management token",
		},
	}
}
