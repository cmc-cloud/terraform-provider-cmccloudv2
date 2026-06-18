package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func createPortAddressPairElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"mac_address": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
func portConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port_id": {
			Type:         schema.TypeString,
			ValidateFunc: validateUUID,
			ForceNew:     true,
			Required:     true,
		},
		"name": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateName,
			Description:  "The name of the port",
		},
		"port_security_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The security of the port",
		},
		"security_group_ids": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "The security group ids",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
		},
		"allowed_address_pairs": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Dashboard allowed address pairs",
			Elem: &schema.Resource{
				Schema: createPortAddressPairElementSchema(),
			},
		},
	}
}
