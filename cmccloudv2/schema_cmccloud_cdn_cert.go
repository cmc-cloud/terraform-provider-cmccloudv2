package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func cdnCertSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cert_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the CDN certificate",
		},
		"cert_data": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Certificate data, generate example from https://en.rakko.tools/tools/46/",
			Sensitive:   true,
		},
		"key_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the private key",
		},
		"key_data": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Private key, generate example from https://en.rakko.tools/tools/46/",
			Sensitive:   true,
		},
		"certificate_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The type of the CDN certificate, `SSL`, `WARP`",
		},
		"common_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The common name of the CDN certificate",
		},
		"expiration_date": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The expiration date of the CDN certificate",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the CDN certificate",
		},
	}
}
