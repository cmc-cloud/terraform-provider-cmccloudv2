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
		},

		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
		},

		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"symmetric", "public", "private", "passphrase", "certificate",
			}, false),
		},

		"content": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
			ForceNew:  true,
			DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
				return strings.TrimSpace(o) == strings.TrimSpace(n)
			},
		},

		"expiration": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"secret_ref": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
