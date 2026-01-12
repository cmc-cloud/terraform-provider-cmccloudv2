package cmccloudv2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func keymanagementsecretSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The ID of the key management container",
		},

		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the key management secret",
		},

		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"symmetric", "public", "private", "passphrase", "certificate",
			}, false),
			Description: "The type of the key management secret",
		},

		"content": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
			ForceNew:  true,
			DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
				return strings.TrimSpace(o) == strings.TrimSpace(n)
			},
			Description: "The content of the key management secret",
		},

		"expiration": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
			Description:  "The expiration time of the key management secret",
		},

		"secret_ref": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The reference of the key management secret",
		},

		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the key management secret",
		},
	}
}
