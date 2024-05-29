package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// $name, $volume_id, $hour, $minute, $interval, $max_keep, $incremental

func databaseAutoBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"schedule_time": {
			Type:     schema.TypeString,
			Required: true,
		},
		"interval": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"max_keep": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"incremental": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_run": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}
