package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// Clone + extend
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
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "",
		},
		"security_group_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the parent security group.",
		},
	}
}

// func securityGroupRuleSchema() map[string]*schema.Schema {
// 	return map[string]*schema.Schema{
// 		"ether_type": {
// 			Type:         schema.TypeString,
// 			Required:     true,
// 			ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, true),
// 			Description:  "Must be IPv4 or IPv6, and addresses represented in CIDR must match the ingress or egress rules.",
// 		},
// 		"direction": {
// 			Type:         schema.TypeString,
// 			Required:     true,
// 			ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, true),
// 		},
// 		"protocol": {
// 			Type:        schema.TypeString,
// 			Required:    true,
// 			Description: "The IP protocol can be represented by a string, an integer, or empty. Valid string or integer values are any or 0, ah or 51, dccp or 33, egp or 8, esp or 50, gre or 47, icmp or 1, icmpv6 or 58, igmp or 2, ipip or 4, ipv6-encap or 41, ipv6-frag or 44, ipv6-icmp or 58, ipv6-nonxt or 59, ipv6-opts or 60, ipv6-route or 43, ospf or 89, pgm or 113, rsvp or 46, sctp or 132, tcp or 6, udp or 17, udplite or 136, vrrp or 112. Additionally, any integer value between [0-255] is also valid. The string any (or integer 0) means all IP protocols",
// 		},
// 		"port_range_max": {
// 			Type:     schema.TypeInt,
// 			Optional: true,
// 			// ValidateFunc: validateAny(
// 			// 	"must either be an empty string, a port number or a valid port range",
// 			// 	validateEmpty,
// 			// 	validatePortNumber,
// 			// 	validatePortRange,
// 			// ),
// 			// Default:     "",
// 			Description: "The maximum port number in the range that is matched by the security group rule. If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be greater than or equal to the port_range_min attribute value. If the protocol is ICMP, this value must be an ICMP code.",
// 		},
// 		"port_range_min": {
// 			Type:     schema.TypeInt,
// 			Optional: true,
// 			// ValidateFunc: validateAny(
// 			// 	"must either be an empty string, a port number or a valid port range",
// 			// 	validateEmpty,
// 			// 	validatePortNumber,
// 			// 	validatePortRange,
// 			// ),
// 			// Default:     "",
// 			Description: "The minimum port number in the range that is matched by the security group rule. If the protocol is TCP, UDP, DCCP, SCTP or UDP-Lite this value must be less than or equal to the port_range_max attribute value. If the protocol is ICMP, this value must be an ICMP type.",
// 		},
// 		"cidr": {
// 			Type:     schema.TypeString,
// 			Optional: true,
// 			ValidateFunc: validateAny(
// 				"must either be an empty string or a valid cidr range",
// 				validateEmpty,
// 				validateIPCidrRange,
// 			),
// 			Default:     "",
// 			Description: "The CIDR that is matched by this security group rule.",
// 		},
// 		"remote_group_id": {
// 			Type:     schema.TypeString,
// 			Optional: true,
// 			Default:  "",
// 			ValidateFunc: validateAny(
// 				"must either be an empty string or a valid UUID",
// 				validateEmpty,
// 				validateUUID,
// 			),
// 			Description: "The securitygroup ID to associate with this securitygroup rule. You can specify either the remote_group_id or cidr attribute in the request ",
// 		},
// 		"description": {
// 			Type:     schema.TypeString,
// 			Optional: true,
// 			Default:  "",
// 		},
// 		"id": {
// 			Type:     schema.TypeString,
// 			Computed: true,
// 			// DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
// 			// 	return true // luôn suppress diff
// 			// },
// 		},
// 	}
// }

func securityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"stateful": {
			Type:     schema.TypeBool,
			Required: true,
		},
		// "rule": {
		// 	Type:     schema.TypeSet,
		// 	Optional: true,
		// 	Computed: true,
		// 	Set:      computeSecGroupV2RuleHash,
		// 	Elem: &schema.Resource{
		// 		Schema: securityGroupRuleSchema(),
		// 	},
		// },
	}
}

// func computeSecGroupV2RuleHash(v interface{}) int {
// 	var buf bytes.Buffer
// 	m := v.(map[string]interface{})

// 	min, ok := m["port_range_min"].(int)
// 	if !ok {
// 		min = 0
// 	}
// 	max, ok := m["port_range_max"].(int)
// 	if !ok {
// 		max = 0
// 	}
// 	if m["cidr"] == nil {
// 		m["cidr"] = ""
// 	}
// 	if m["remote_group_id"] == nil {
// 		m["remote_group_id"] = ""
// 	}

// 	buf.WriteString(fmt.Sprintf("%s-", m["direction"].(string)))
// 	buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))
// 	buf.WriteString(fmt.Sprintf("%d-", min))
// 	buf.WriteString(fmt.Sprintf("%d-", max))
// 	buf.WriteString(fmt.Sprintf("%s-", m["cidr"].(string)))
// 	buf.WriteString(fmt.Sprintf("%s-", m["remote_group_id"].(string)))
// 	buf.WriteString(fmt.Sprintf("%s-", m["ether_type"].(string)))

// 	hash := int(crc32.ChecksumIEEE(buf.Bytes()))

// 	// In ra log hoặc console (nếu bạn có logger tốt hơn thì thay thế)
// 	gocmcapiv2.Logs(fmt.Sprintf("[DEBUG] Rule Hash Input: %s -> Hash: %d", buf.String(), hash))

// 	return hash
// }
