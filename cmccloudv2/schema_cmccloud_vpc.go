package cmccloudv2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func vpcSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the VPC",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the VPC",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
			Description:  "Name of billing mode, monthly or hourly",
		},
		"tags": tagSchema(),
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
			ValidateFunc: validateIPCidrRange,
			Description:  "CIDR representing IP range for this VPC",
		},
	}
}
