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
		},
		"message": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"detection": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"BLOCK", "LOG"}, false),
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"match_request_body": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_get_arguments": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_http_headers": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_filename": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_url": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_name_check": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_header_cookie": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_header_content_type": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_header_user_agent": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_header_accept_encoding": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"match_header_connection": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}
