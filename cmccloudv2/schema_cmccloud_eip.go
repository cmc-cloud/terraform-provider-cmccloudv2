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
			Description:  "The billing mode of the EIP. A valid value is monthly or hourly.",
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
		"tags": tagSchema(),
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the EIP. Changing this updates the EIP's description.",
		},
		"dns_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The DNS domain of the EIP.",
		},
		"dns_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The DNS name of the EIP.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the EIP.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the EIP.",
		},
		"eip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The IP address of the EIP after it is allocated.",
		},
	}
}
