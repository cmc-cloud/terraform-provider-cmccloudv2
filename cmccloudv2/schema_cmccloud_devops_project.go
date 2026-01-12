package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func devopsProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Devops Project",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The description of the Devops Project",
		},
		"is_default": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "If true, the Devops Project is the default one",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the Devops Project",
		},
	}
}
