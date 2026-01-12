package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func containerRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"devops_project_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the DevOps project",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the container registry",
		},
		"uri": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The URI of the container registry",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the container registry",
		},
	}
}
