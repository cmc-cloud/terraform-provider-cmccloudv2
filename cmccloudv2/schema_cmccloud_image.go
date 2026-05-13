package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func imageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the image",
		},
		"volume_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The id of volume that you want to create from",
		},
		"disk_format": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "raw",
			ValidateFunc: validation.StringInSlice([]string{"qcow2", "raw", "vdi", "vhd", "vhdx", "vmdk"}, false),
			Description:  "The type of disk format",
		},
		"force": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
			Description: "Set true if force create image",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "OpenSearch status",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time",
		},
	}
}
