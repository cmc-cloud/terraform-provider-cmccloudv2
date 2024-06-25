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
		},
		"ip": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateIPAddress,
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ignore", "deny"}, false),
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
