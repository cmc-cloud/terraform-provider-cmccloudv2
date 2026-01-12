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
			Description:  "The ID of the Server to attach the Interface to",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the Subnet to attach the Interface to",
		},
		"ip_address": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateIPAddress,
			Description:  "The IP address of the Interface. If not provided, a random IP address will be assigned",
		},
	}
}
