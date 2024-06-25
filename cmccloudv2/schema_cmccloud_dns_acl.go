package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dnsAclSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			// ValidateFunc: validateUUID,
		},
		"record_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "MX", "NS", "TXT", "CNAME"}, true),
		},
		"source_ip": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.Any(validateIPAddress, validateIPCidrRange),
		},
		"domain": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateDomainName,
		},
		"action": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "none",
			ValidateFunc: validation.StringInSlice([]string{"block", "allow"}, false),
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
