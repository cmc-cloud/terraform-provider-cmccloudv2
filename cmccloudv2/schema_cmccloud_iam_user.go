package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func iamUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUserName,
			ForceNew:     true,
			Description:  "The username of the IAM user.",
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
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The password for the IAM user.",
		},
		"email": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateEmail,
			Description:  "The email address of the IAM user.",
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
