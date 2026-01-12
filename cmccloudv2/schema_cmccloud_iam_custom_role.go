package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func iamCustomRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the custom role",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the custom role",
		},
		"content": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The json content of custom role policy.",
			ValidateFunc: validation.StringIsJSON,
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the custom role",
		},
	}
}
