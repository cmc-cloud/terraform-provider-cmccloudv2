package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func iamCustomRoleAssignmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the project.",
		},
		"group_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The id of the group.",
		},
		"custom_role_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the custom role.",
		},
	}
}
