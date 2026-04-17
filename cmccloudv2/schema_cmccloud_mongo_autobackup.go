package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// $name, $volume_id, $hour, $minute, $interval, $max_keep, $incremental

func mongoAutoBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the mongo instance",
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
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the auto backup",
		},
		"next_run": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The last run time of the auto backup",
		},
	}
}
