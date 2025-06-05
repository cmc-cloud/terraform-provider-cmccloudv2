package cmccloudv2

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func lbL7policySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"action": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"REDIRECT_TO_POOL", "REDIRECT_TO_URL", "REJECT", "REDIRECT_PREFIX",
			}, true),
		},

		"listener_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"position": {
			Type:     schema.TypeInt,
			Optional: true,
			// Computed: true,
		},

		"redirect_prefix": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"redirect_url", "redirect_pool_id"},
			Optional:      true,
		},

		"redirect_pool_id": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"redirect_url", "redirect_prefix"},
			Optional:      true,
		},

		"redirect_url": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"redirect_pool_id", "redirect_prefix"},
			Optional:      true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(string)
				_, err := url.ParseRequestURI(value)
				if err != nil {
					errors = append(errors, fmt.Errorf("URL is not valid: %s", err))
				}
				return
			},
		},

		"redirect_http_code": {
			Type:          schema.TypeInt,
			ConflictsWith: []string{"redirect_pool_id"},
			Optional:      true,
			// Computed:      true,
			ValidateFunc: validation.IntInSlice([]int{301, 302, 303, 307, 308}),
		},
	}
}
