package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// $name, $volume_id, $hour, $minute, $interval, $max_keep, $incremental

func volumeAutoBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the auto backup",
		},
		"volume_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the Volume to backup",
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
		"incremental": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "If true, it is incremental backup, if false, it is full backup. Default is true",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the auto backup",
		},
		"last_run": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The last run time of the auto backup",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The size of the volume in GB",
		},
	}
}
