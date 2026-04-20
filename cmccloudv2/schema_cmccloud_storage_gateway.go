package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func storageGatewaySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "Name of the storage gateway",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the storage gateway",
		},
		"protocol_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"nfs"}, false),
			Description:  "Protocol type of the storage gateway. Only `nfs` is currently supported",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
			Description:  "ID of subnet",
		},
		"bucket": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The name of the S3 bucket used as the backend storage for the gateway. All files written to the shared path (via NFS) will be uploaded to this bucket",
		},
		"tags": tagSchema(),
		"command_line": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The command or mount instruction used by clients to access the shared storage from the gateway",
		},
		"shared_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The local directory path on the storage gateway that will be shared and mapped to S3",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the storage gateway",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the storage gateway",
		},
	}
}
