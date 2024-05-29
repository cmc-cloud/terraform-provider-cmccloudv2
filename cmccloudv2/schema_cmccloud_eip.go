package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func eipSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode, // validation.StringInSlice([]string{"monthly", "hourly"}, true)
			Default:      "monthly",
			Optional:     true,
		},
		"domestic_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Domestic bandwidth in Mbps",
		},
		"inter_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "International bandwidth in Mbps",
		},
		// "port_forwardings": {
		// 	Type:     schema.TypeList,
		// 	Optional: true,
		// 	Elem: &schema.Resource{
		// 		Schema: createEipPortForwardingsElementSchema(),
		// 	},
		// },
		"tags": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns_domain": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"eip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
