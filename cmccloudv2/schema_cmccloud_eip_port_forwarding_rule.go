package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createEipPortForwardingRuleElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"eip_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "sctp", "icmp", "icmp6", "dccp"}, true),
		},
		"internal_ip_address": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateIPAddress,
		},
		"internal_port_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
		},
		"internal_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"external_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"internal_port_range": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"external_port_range": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
