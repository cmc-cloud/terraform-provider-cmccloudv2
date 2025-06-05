package cmccloudv2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func lbL7policyRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"COOKIE", "FILE_TYPE", "HEADER", "HOST_NAME",
				"PATH",
				//"SSL_CONN_HAS_CERT", "SSL_VERIFY_RESULT",
				//"SSL_DN_FIELD",
			}, true),
		},
		"compare_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"CONTAINS", "STARTS_WITH", "ENDS_WITH", "EQUAL_TO", "REGEX",
			}, true),
		},
		"l7policy_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"provisioning_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"value": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				if len(v.(string)) == 0 {
					errors = append(errors, fmt.Errorf("'value' field should not be empty"))
				}
				return
			},
		},

		"key": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"invert": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
		},
	}
}
