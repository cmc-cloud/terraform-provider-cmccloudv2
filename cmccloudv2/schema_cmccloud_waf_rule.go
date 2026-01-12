package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func wafruleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"waf_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the WAF",
		},
		"message": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The message of the WAF rule",
		},
		"detection": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The detection of the WAF rule",
		},
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"BLOCK", "LOG"}, false),
			Description:  "The action of the WAF rule, BLOCK or LOG",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The description of the WAF rule",
		},
		"match_request_body": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the request body",
		},
		"match_get_arguments": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the GET arguments",
		},
		"match_http_headers": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the HTTP headers",
		},
		"match_filename": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the filename",
		},
		"match_url": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the URL",
		},
		"match_name_check": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the name check",
		},
		"match_header_cookie": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the header cookie",
		},
		"match_header_content_type": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the header content type",
		},
		"match_header_user_agent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the header user agent",
		},
		"match_header_accept_encoding": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the header accept encoding",
		},
		"match_header_connection": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If true, it will match the header connection",
		},
	}
}
