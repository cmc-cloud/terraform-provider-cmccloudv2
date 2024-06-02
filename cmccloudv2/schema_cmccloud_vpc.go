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
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
		},
		"tags": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
			ValidateFunc: validateIPCidrRange,
		},
	}
}
