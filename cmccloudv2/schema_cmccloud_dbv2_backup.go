package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dbv2BackupSchema(dbType string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: dbType + " instance id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of the backup",
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
