package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func vipBindingSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vip_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the subnet to attach the VIP to",
		},
		"port_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the port to attach the VIP to",
		},
	}
}
