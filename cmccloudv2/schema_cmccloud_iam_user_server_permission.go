package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func iamUserServerPermissionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The username of the user",
		},
		"server_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the server",
		},
		"blocked": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the user is blocked from the server",
		},
		"allow_view": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the user is allowed to view the server",
		},
		"allow_edit": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the user is allowed to edit the server",
		},
		"allow_create": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the user is allowed to create resources on the server",
		},
		"allow_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the user is allowed to delete resources on the server",
		},
	}
}
