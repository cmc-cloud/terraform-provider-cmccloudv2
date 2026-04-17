package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func postgresBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Postgres instance id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Postgres backup name",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Size in bytes of backup",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Type of backup, full or incremental",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time",
		},
	}
}
