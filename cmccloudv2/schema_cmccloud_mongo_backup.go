package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func mongoBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mongo instance id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mongo backup name",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Size in bytes of backup",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time",
		},
	}
}
