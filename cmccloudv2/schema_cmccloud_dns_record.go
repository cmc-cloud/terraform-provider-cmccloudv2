package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dnsRecordSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			// ValidateFunc: validateUUID,
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"ttl": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"load_balance_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "none",
			ValidateFunc: validation.StringInSlice([]string{"weighted", "none"}, false),
		},
		"ips": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: createDnsRecordIpSchema(),
			},
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
