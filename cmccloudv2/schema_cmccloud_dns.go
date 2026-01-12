package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dnsZoneSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateDomainName,
			ForceNew:     true,
			Description:  "The domain of the DNS zone",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"primary"}, false),
			Description:  "The type of the DNS zone. Currently, only `primary` type is supported",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the DNS zone",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the DNS zone",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The updated time of the DNS zone",
		},
	}
}
