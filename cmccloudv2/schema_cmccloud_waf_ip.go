package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func wafipSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"waf_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the WAF",
		},
		"ip": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateIPAddress,
			Description:  "The IP address to add to the WAF",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ignore", "deny"}, false),
			Description:  "The type of the IP address, ignore or deny",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the IP address",
		},
	}
}
