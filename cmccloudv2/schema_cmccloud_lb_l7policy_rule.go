package cmccloudv2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func lbL7policyRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"l7policy_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the L7 policy to attach the L7 policy rule to",
		},
		"type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"COOKIE", "FILE_TYPE", "HEADER", "HOST_NAME",
				"PATH",
				//"SSL_CONN_HAS_CERT", "SSL_VERIFY_RESULT",
				//"SSL_DN_FIELD",
			}, true),
			Description: "The type of the L7 policy rule",
		},
		"compare_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"CONTAINS", "STARTS_WITH", "ENDS_WITH", "EQUAL_TO", "REGEX",
			}, true),
			Description: "The compare type of the L7 policy rule",
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
			Description: "The value of the L7 policy rule",
		},

		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The key of the L7 policy rule",
		},

		"invert": {
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "The invert of the L7 policy rule",
		},

		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the L7 policy rule",
		},

		"provisioning_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The provisioning status of the L7 policy rule",
		},
		"operating_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The operating status of the L7 policy rule",
		},
	}
}
