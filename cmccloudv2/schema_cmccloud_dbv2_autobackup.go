package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dbv2AutoBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the database instance",
		},
		"schedule_time": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The schedule time in HH:mm format, eg: 19:05",
		},
		"interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The interval of the auto backup",
		},
		"max_keep": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "The maximum number of backups to keep",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the auto backup",
		},
		"next_run": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The next run time of the auto backup",
		},
	}
}
