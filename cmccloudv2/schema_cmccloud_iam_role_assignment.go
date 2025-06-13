package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func iamRoleAssignmentSchema() map[string]*schema.Schema {
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
		// "group_name": {
		// 	Type:        schema.TypeString,
		// 	Required:    true,
		// 	ForceNew:    true,
		// 	Description: "The name of the group.",
		// },
		"role_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the role.",
		},
	}
}
