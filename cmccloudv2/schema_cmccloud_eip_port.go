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
			Description:  "The ID of the EIP to attach the port to",
		},
		"port_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the port to attach the EIP to",
		},
		"fix_ip_address": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateIPAddress,
			Description:  "The IP address of the port to attach the EIP to",
		},
	}
}
