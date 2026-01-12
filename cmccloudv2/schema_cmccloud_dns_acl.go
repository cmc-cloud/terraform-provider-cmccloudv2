package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dnsAclSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the DNS zone",
		},
		"record_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "MX", "NS", "TXT", "CNAME"}, true),
			Description:  "The type of the DNS record",
		},
		"source_ip": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.Any(validateIPAddress, validateIPCidrRange),
			Description:  "The source IP of the DNS ACL",
		},
		"domain": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateDomainName,
			Description:  "The domain of the DNS ACL",
		},
		"action": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "none",
			ValidateFunc: validation.StringInSlice([]string{"block", "allow"}, false),
			Description:  "The action of the DNS ACL",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the DNS ACL",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The updated time of the DNS ACL",
		},
	}
}
