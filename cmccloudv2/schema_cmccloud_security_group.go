package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func securityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the security group",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The description of the security group",
		},
		"stateful": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The stateful of the security group, if true, the security group is stateful",
		},
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
