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
			Description:  "The ID of the EIP to attach the port forwarding rule to",
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "sctp", "icmp", "icmp6", "dccp"}, true),
			Description:  "The protocol for the port forwarding rule. A valid value is tcp, udp, sctp, icmp, icmp6, or dccp.",
		},
		"internal_ip_address": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateIPAddress,
			Description:  "The fixed IPv4 address of the port associated to the EIP port forwarding rule.",
		},
		"internal_port_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the internal port to attach the port forwarding rule to",
		},
		"internal_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "The TCP/UDP/other protocol port number of the port fixed IP address associated to the floating ip port forwarding.. A valid value is between 1 and 65535.",
		},
		"external_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "The external port number of the port forwarding rule. A valid value is between 1 and 65535.",
		},
		"internal_port_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The internal port range of the port forwarding rule. A valid value is a comma-separated list of ports or a range of ports (e.g., 1000-2000).",
		},
		"external_port_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The external port range of the port forwarding rule. A valid value is a comma-separated list of ports or a range of ports (e.g., 1000-2000).",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the port forwarding rule. Changing this updates the port forwarding rule's description.",
		},
	}
}
