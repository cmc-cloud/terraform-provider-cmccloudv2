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
			ValidateFunc: validation.NoZeroValues,
		},

		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},

		"expiration": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"type": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				"symmetric", "public", "private", "passphrase", "certificate", "opaque",
			}, false),
		},

		"content": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			ForceNew:  true,
			DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
				return strings.TrimSpace(o) == strings.TrimSpace(n)
			},
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
