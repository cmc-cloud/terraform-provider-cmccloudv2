package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func serverInterfaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"server_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"ip_address": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateIPAddress,
		},
	}
}
