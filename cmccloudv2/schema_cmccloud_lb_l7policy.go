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
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the L7 policy",
		},
		"action": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"REDIRECT_TO_POOL", "REDIRECT_TO_URL", "REJECT", "REDIRECT_PREFIX",
			}, true),
			Description: "The action of the L7 policy",
		},

		"listener_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the listener to attach the L7 policy to",
		},

		"position": {
			Type:     schema.TypeInt,
			Optional: true,
			// Computed: true,
			Description: "The position of this policy on the listener. Positions start at 1.",
		},

		"redirect_prefix": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"redirect_url", "redirect_pool_id"},
			Optional:      true,
			Description:   "Requests matching this policy will be redirected to this Prefix URL. Only valid if action is REDIRECT_PREFIX.",
		},

		"redirect_pool_id": {
			Type:          schema.TypeString,
			ConflictsWith: []string{"redirect_url", "redirect_prefix"},
			Optional:      true,
			Description:   "Requests matching this policy will be redirected to the pool with this ID. Only valid if action is REDIRECT_TO_POOL. The pool has some restrictions",
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
			Description: "Requests matching this policy will be redirected to this URL. Only valid if action is REDIRECT_TO_URL.",
		},

		"redirect_http_code": {
			Type:          schema.TypeInt,
			ConflictsWith: []string{"redirect_pool_id"},
			Optional:      true,
			// Computed:      true,
			ValidateFunc: validation.IntInSlice([]int{301, 302, 303, 307, 308}),
			Description:  "Requests matching this policy will be redirected to the specified URL or Prefix URL with the HTTP response code. Valid if action is REDIRECT_TO_URL or REDIRECT_PREFIX. Valid options are: 301, 302, 303, 307, or 308.",
		},
	}
}
