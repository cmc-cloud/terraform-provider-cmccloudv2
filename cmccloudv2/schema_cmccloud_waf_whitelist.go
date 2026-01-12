package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func wafwhitelistSchema() map[string]*schema.Schema {
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
			Description: "The message of the WAF whitelist",
		},
		"detection": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The detection of the WAF whitelist",
		},
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"BLOCK", "LOG"}, false),
			Description:  "The action of the WAF whitelist, BLOCK or LOG",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The description of the WAF whitelist",
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
		"match_header_var": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"Cookie", "Content-Type", "User-Agent", "Accept-Encoding", "Connection"}, false),
			Description:  "The header variable to match, Cookie, Content-Type, User-Agent, Accept-Encoding, Connection",
		},
	}
}
