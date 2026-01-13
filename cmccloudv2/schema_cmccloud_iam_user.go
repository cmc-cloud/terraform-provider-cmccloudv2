package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func iamUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"short_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUserName,
			ForceNew:     true,
			Description:  "The short name of the IAM user. The full name of the IAM user is `domain_name` + `_` + `short_name`, example: `3hr4enae2tvg_dev`",
		},
		"first_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The first name of the IAM user.",
		},
		"last_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The last name of the IAM user.",
		},
		"password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validatePassword,
			Description:  "The password for the IAM user. Minimum Length 12, Require at least one uppercase character, one lowercase character, one number, one special character",
		},
		"email": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateEmail,
			Description:  "The email address of the IAM user.",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "The status of the IAM user",
		},
		"username": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The username of the IAM user",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the IAM user",
		},
	}
}
