package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func eipportSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"eip_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"port_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"fix_ip_address": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateIPAddress,
		},
	}
}
