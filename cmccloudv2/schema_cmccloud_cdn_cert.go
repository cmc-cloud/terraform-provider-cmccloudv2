package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func cdnCertSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cert_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"cert_data": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Certificate data, generate example from https://en.rakko.tools/tools/46/",
			Sensitive:   true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"key_data": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Private key, generate example from https://en.rakko.tools/tools/46/",
			Sensitive:   true,
		},
		"certificate_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"common_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"expiration_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
