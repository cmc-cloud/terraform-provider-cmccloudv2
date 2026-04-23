package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func apigatewayBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The id of api gateway instance",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the backup",
		},
		"size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "The size of the backup in byte",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the backup",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the backup",
		},
	}
}
