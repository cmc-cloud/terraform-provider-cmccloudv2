package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func wafcertSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"cert_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cert_data": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Certificate data, generate example from https://en.rakko.tools/tools/46/",
			Sensitive:   true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"key_data": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Private key, generate example from https://en.rakko.tools/tools/46/",
			ForceNew:    true,
			Sensitive:   true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}
