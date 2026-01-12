package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func volumeSnapshotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the snapshot",
		},
		"volume_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the Volume to snapshot",
		},
		"force": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "If true, it will force the snapshot to be created even if the volume is not in a stable state. Default is false",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the snapshot",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the snapshot",
		},
		"volume_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of the volume",
		},
		"volume_size": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "The size of the snapshoted volume in GB",
		},
		"real_size_gb": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The real size of the snapshot in GB",
		},
	}
}
