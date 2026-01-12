package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func volumeBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the backup",
		},
		"volume_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the Volume to backup",
		},
		"incremental": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "If true, it is incremental backup, if false, it is full backup. Default is false",
		},
		"force": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "If true, it will force the backup to be created even if the volume is not in a stable state. Default is false",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the backup",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the backup",
		},
		"volume_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the volume",
		},
		"volume_size": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "The size of backuped volume in GB",
		},
		"real_size_gb": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "The real size of the backup in GB",
		},
	}
}
