package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dnsRecordSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the DNS zone",
		},
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The type of the DNS record",
		},
		"domain": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The domain of the DNS record",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The TTL of the DNS record",
		},
		"load_balance_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "none",
			ValidateFunc: validation.StringInSlice([]string{"weighted", "none"}, false),
			Description:  "The load balance type of the DNS record",
		},
		"ips": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: createDnsRecordIpSchema(),
			},
			Description: "The IPs of the DNS record",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the DNS record",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The updated time of the DNS record",
		},
	}
}

func createDnsRecordIpSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"weight": {
			Type:     schema.TypeInt,
			Optional: true,
			// Default:  1,
			ForceNew: true,
		},
	}
}
