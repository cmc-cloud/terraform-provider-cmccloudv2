package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func securityGroupRuleWithGroupIdSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ether_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, true),
			Description:  "Must be IPv4 or IPv6, and addresses represented in CIDR must match the ingress or egress rules.",
		},
		"direction": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, true),
			Description:  "The direction of the security group rule. Must be ingress (inbound) or egress (outbound)",
		},
		"protocol": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The IP protocol can be represented by a string, an integer, or empty. Valid string or integer values are any or 0, ah or 51, dccp or 33, egp or 8, esp or 50, gre or 47, icmp or 1, icmpv6 or 58, igmp or 2, ipip or 4, ipv6-encap or 41, ipv6-frag or 44, ipv6-icmp or 58, ipv6-nonxt or 59, ipv6-opts or 60, ipv6-route or 43, ospf or 89, pgm or 113, rsvp or 46, sctp or 132, tcp or 6, udp or 17, udplite or 136, vrrp or 112. Additionally, any integer value between [0-255] is also valid. The string any (or integer 0) means all IP protocols",
		},
		"port_range_max": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "The maximum port number in the range that is matched by the security group rule. If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be greater than or equal to the port_range_min attribute value. If the protocol is ICMP, this value must be an ICMP code.",
		},
		"port_range_min": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "The minimum port number in the range that is matched by the security group rule. If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be less than or equal to the port_range_max attribute value. If the protocol is ICMP, this value must be an ICMP type.",
		},
		"cidr": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validateAny(
				"must either be an empty string or a valid cidr range",
				validateEmpty,
				validateIPCidrRange,
			),
			Default:     "",
			Description: "The CIDR that is matched by this security group rule.",
		},
		"remote_group_id": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			ForceNew: true,
			ValidateFunc: validateAny(
				"must either be an empty string or a valid UUID",
				validateEmpty,
				validateUUID,
			),
			Description: "The securitygroup ID to associate with this securitygroup rule. You can specify either the remote_group_id or cidr attribute in the request ",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "",
			Description: "A description of the rule. Changing this creates a new security group rule.",
		},
		"security_group_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the parent security group.",
		},
	}
}
